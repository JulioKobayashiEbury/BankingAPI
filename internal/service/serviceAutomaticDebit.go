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

type ServiceAutoDebit interface {
	ProcessNewAutomaticDebit(autoDebit *automaticdebit.AutomaticDebitRequest) (*automaticdebit.AutomaticDebitResponse, *model.Erro)
	CheckAutomaticDebits()
}

type serviceAutoDebitImpl struct {
	autoDebitFirestore model.RepositoryInterface
	withdrawalService  ServiceWithdrawal
}

func NewAutoDebitImpl(autodebitDB model.RepositoryInterface, withdrawal ServiceWithdrawal) ServiceAutoDebit {
	return serviceAutoDebitImpl{
		autoDebitFirestore: autodebitDB,
		withdrawalService:  withdrawal,
	}
}

func (debitService serviceAutoDebitImpl) ProcessNewAutomaticDebit(autoDebit *automaticdebit.AutomaticDebitRequest) (*automaticdebit.AutomaticDebitResponse, *model.Erro) {
	if !isValidDate(autoDebit.Expiration_date) {
		log.Warn().Msg("Invalid date format")
		return nil, &model.Erro{Err: errors.New("Invalid date format"), HttpCode: http.StatusBadRequest}
	}
	responseID, err := debitService.autoDebitFirestore.Create(autoDebit)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Automatic debit created: " + *responseID)
	obj, err := debitService.autoDebitFirestore.Get(responseID)
	if err != nil {
		return nil, err
	}
	autoDebitResponse, ok := obj.(automaticdebit.AutomaticDebitResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return &autoDebitResponse, nil
}

func isValidDate(date string) bool {
	fmt.Println(date)
	_, err := time.Parse(model.TimeLayout, date)
	return err == nil
}

func (autoDebitService serviceAutoDebitImpl) CheckAutomaticDebits() {
	log.Info().Msg("Checking for auto debits...")
	obj, err := autoDebitService.autoDebitFirestore.GetAll()
	if err != nil {
		log.Error().Msg(err.Err.Error())
		return
	}
	autoDebitList, ok := obj.(*[]*automaticdebit.AutomaticDebitResponse)
	if !ok {
		log.Error().Msg("Error getting automatic debit list as data type returned is wrong")
		return
	}
	for index := 0; index < len(*(autoDebitList)); index++ {
		autoDebit := (*autoDebitList)[index]
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

				if err := autoDebitService.autoDebitFirestore.Update(&autoDebit.Debit_id); err != nil {
					return
				}
				log.Warn().Msg("Debit is expired")
				return
			} else {
				if autoDebit.Debit_day == uint16(time.Now().Day()) {
					_, err := autoDebitService.withdrawalService.ProcessWithdrawal(&withdrawal.WithdrawalRequest{
						Account_id: autoDebit.Account_id,
						Client_id:  autoDebit.Client_id,
						Agency_id:  autoDebit.Agency_id,
						Withdrawal: autoDebit.Value,
					})
					if err != nil {
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
