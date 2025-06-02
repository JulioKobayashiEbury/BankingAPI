package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	model "BankingAPI/internal/model"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/withdrawal"

	"github.com/rs/zerolog/log"
)

const (
	//"2006-01-02T15:04:05+07:00"
	timeLayout = time.RFC3339
)

func ProcessNewAutomaticDebit(autoDebit *automaticdebit.AutomaticDebit) *model.Erro {
	// map new debit
	if !IsValidDate(autoDebit.Expiration_date) {
		log.Warn().Msg("Invalid date format")
		return &model.Erro{Err: errors.New("Invalid date format"), HttpCode: http.StatusBadRequest}
	}
	// insert in db
	database := &automaticdebit.AutoDebitFirestore{
		AutoDebit: autoDebit,
	}
	if err := database.Create(); err != nil {
		return err
	}
	log.Info().Msg("Automatic debit created: " + database.AutoDebit.Debit_id)
	return nil
}

func IsValidDate(date string) bool {
	fmt.Println(date)
	_, err := time.Parse(model.TimeLayout, date)
	return err == nil
}

func CheckAutomaticDebits() {
	log.Info().Msg("Checking for auto debits...")
	database := &automaticdebit.AutoDebitFirestore{}
	if err := database.GetAll(); err != nil {
		log.Error().Msg(err.Err.Error())
		return
	}
	for index := 0; index < len(*(database.Slice)); index++ {
		autoDebit := (*database.Slice)[index]
		if !autoDebit.Status {
			log.Warn().Msg("Debit is expired")
			return
		} else {
			expirationDate, err := time.Parse(timeLayout, autoDebit.Expiration_date)
			if err != nil {
				log.Error().Msg(err.Error())
				return
			}
			if expirationDate.Unix() > time.Now().Unix() {
				database.AddUpdate("status", false)
				if err := database.Update(); err != nil {
					return
				}
				log.Warn().Msg("Debit is expired")
				return
			} else {
				if autoDebit.Debit_day == uint16(time.Now().Day()) {
					if err := ProcessWithdrawal(&withdrawal.WithdrawalRequest{
						Account_id: autoDebit.Account_id,
						Client_id:  autoDebit.Client_id,
						Agency_id:  autoDebit.Agency_id,
						Withdrawal: autoDebit.Value,
					}); err != nil {
						log.Error().Msg(err.Err.Error())
						return
					}
					log.Info().Msg("Auto debit is logged as Withdrawal together with Account")
				}
			}
		}

	}
}

func moveToExpiredDebits() {
}
