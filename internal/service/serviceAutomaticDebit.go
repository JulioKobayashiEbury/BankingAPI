package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

const (
	//"2006-01-02T15:04:05+07:00"
	timeLayout = time.RFC3339
)

func ProcessNewAutomaticDebit(autoDebit *model.AutomaticDebitRequest) *model.Erro {
	// map new debit
	if !IsValidDate(autoDebit.Expiration_date) {
		log.Warn().Msg("Invalid date format")
		return &model.Erro{Err: errors.New("Invalid date format"), HttpCode: http.StatusBadRequest}
	}

	debitMap := map[string]interface{}{
		"account_id":      autoDebit.Account_id,
		"client_id":       autoDebit.Client_id,
		"agency_id":       autoDebit.Agency_id,
		"value":           autoDebit.Value,
		"debit_day":       autoDebit.Debit_day,
		"expiration_date": autoDebit.Expiration_date,
		"register_date":   time.Now().Format(timeLayout),
	}
	// insert in db
	var debitId string
	if err := repository.CreateObject(&debitMap, repository.AutoDebit, &debitId); err != nil {
		return err
	}
	log.Info().Msg("Automatic debit created: " + debitId)
	return nil
}

func IsValidDate(date string) bool {
	fmt.Println(date)
	_, err := time.Parse(timeLayout, date)
	return err == nil
}

func CheckAutomaticDebits() {
	log.Info().Msg("Checking for auto debits...")
	docsSnapshots, err := repository.GetAllByTypeDB(repository.AutoDebit)
	if err != nil {
		log.Error().Msg(err.Err.Error())
		return
	}
	for index := 0; index < len(docsSnapshots); index++ {
		docSnap := docsSnapshots[index]
		log.Info().Msg("Auto debit found: " + docSnap.Ref.ID)
		var autoDebit model.AutomaticDebitResponse
		if err := docSnap.DataTo(&autoDebit); err != nil {
			log.Warn().Msg(fmt.Sprintf("Failed to unmarshal document %s: %v", docSnap.Ref.ID, err))
			return
		}
		expirationDate, err := time.Parse(timeLayout, autoDebit.Expiration_date)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		if expirationDate.Unix() > time.Now().Unix() {
			log.Warn().Msg("Debit is expired")
			return
		}
		if autoDebit.Debit_day == uint16(time.Now().Day()) {
			newBalance, err := ProcessWithdrawal(&model.WithdrawalRequest{
				Account_id: autoDebit.Account_id,
				Client_id:  autoDebit.Client_id,
				Agency_iD:  autoDebit.Agency_id,
				Withdrawal: autoDebit.Value,
			})
			if err != nil {
				log.Error().Msg(err.Err.Error())
				return
			}
			logDebitWithdrawal := map[string]interface{}{
				"debit_id":        autoDebit.Debit_id,
				"account_id":      autoDebit.Account_id,
				"client_id":       autoDebit.Client_id,
				"agency_id":       autoDebit.Agency_id,
				"value":           autoDebit.Value,
				"debit_day":       time.Now().Format(timeLayout),
				"expiration_date": autoDebit.Expiration_date,
				"balance":         newBalance,
			}
			var logDebitID string
			if err := repository.CreateObject(&logDebitWithdrawal, repository.AutoDebitLog, &logDebitID); err != nil {
				log.Error().Msg(err.Err.Error())
				return
			}
			log.Info().Msg("Auto debit is logged: " + logDebitID)
		}
	}
}

func moveToExpiredDebits() {
}
