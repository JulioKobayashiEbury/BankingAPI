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

type serviceAutoDebitImpl struct {
	autoDebitFirestore model.RepositoryInterface
	withdrawalService  WithdrawalService
}

func NewAutoDebit(autodebitDB model.RepositoryInterface, withdrawal WithdrawalService) AutomaticDebitService {
	return serviceAutoDebitImpl{
		autoDebitFirestore: autodebitDB,
		withdrawalService:  withdrawal,
	}
}

func (service serviceAutoDebitImpl) Create(autodebitRequest *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro) {
	obj, err := service.autoDebitFirestore.Create(autodebitRequest)
	if err != nil {
		return nil, err
	}
	automaticDebit, ok := obj.(*automaticdebit.AutomaticDebit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return automaticDebit, nil
}

func (service serviceAutoDebitImpl) Delete(id *string) *model.Erro {
	if err := service.autoDebitFirestore.Delete(id); err != nil {
		return err
	}
	return nil
}

func (service serviceAutoDebitImpl) Get(id *string) (*automaticdebit.AutomaticDebit, *model.Erro) {
	obj, err := service.autoDebitFirestore.Get(id)
	if err != nil {
		return nil, err
	}
	autodebit, ok := obj.(*automaticdebit.AutomaticDebit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return autodebit, nil
}

func (service serviceAutoDebitImpl) GetAll() (*[]automaticdebit.AutomaticDebit, *model.Erro) {
	obj, err := service.autoDebitFirestore.GetAll()
	if err != nil {
		return nil, err
	}
	autodebits, ok := obj.(*[]automaticdebit.AutomaticDebit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return autodebits, nil
}

func (service serviceAutoDebitImpl) Update(autodebitRequest *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro) {
	if err := service.autoDebitFirestore.Update(autodebitRequest); err != nil {
		return nil, err
	}
	return autodebitRequest, nil
}

func (service serviceAutoDebitImpl) Status(id *string, status bool) *model.Erro {
	autodebit, err := service.Get(id)
	if err != nil {
		return nil
	}
	autodebit.Status = status
	if err := service.autoDebitFirestore.Update(autodebit); err != nil {
		return err
	}
	return nil
}

func (service serviceAutoDebitImpl) ProcessNewAutomaticDebit(autoDebit *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro) {
	if !isValidDate(autoDebit.Expiration_date) {
		log.Warn().Msg("Invalid date format")
		return nil, &model.Erro{Err: errors.New("Invalid date format"), HttpCode: http.StatusBadRequest}
	}
	obj, err := service.autoDebitFirestore.Create(autoDebit)
	if err != nil {
		return nil, err
	}
	autodebitResponse, ok := obj.(*automaticdebit.AutomaticDebit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Automatic debit created: " + autodebitResponse.Debit_id)
	return autodebitResponse, nil
}

func isValidDate(date string) bool {
	fmt.Println(date)
	_, err := time.Parse(model.TimeLayout, date)
	return err == nil
}

func (service serviceAutoDebitImpl) CheckAutomaticDebits() {
	log.Info().Msg("Checking for auto debits...")
	autoDebitList, err := service.GetAll()
	if err != nil {
		log.Error().Msg(err.Err.Error())
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
			if expirationDate.Unix() < time.Now().Unix() {
				if err := service.Status(&autoDebit.Debit_id, false); err != nil {
					log.Error().Msg("Failed to update automatic debit status to expired")
					return
				}
				log.Warn().Msg("Debit is expired")
				return
			} else {
				if autoDebit.Debit_day == uint16(time.Now().Day()) {
					_, err := service.withdrawalService.ProcessWithdrawal(&withdrawal.Withdrawal{
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
