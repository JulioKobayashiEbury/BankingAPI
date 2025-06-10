package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/deposit"

	"github.com/rs/zerolog/log"
)

type depositImpl struct {
	depositDatabase model.RepositoryInterface
	accountService  AccountService
}

func NewDepositService(depositDB model.RepositoryInterface, accountServe AccountService) DepositService {
	return depositImpl{
		depositDatabase: depositDB,
		accountService:  accountServe,
	}
}

func (service depositImpl) Create(depositRequest *deposit.Deposit) (*deposit.Deposit, *model.Erro) {
	obj, err := service.depositDatabase.Create(depositRequest)
	if err != nil {
		return nil, err
	}
	depositResponse, ok := obj.(*deposit.Deposit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return depositResponse, nil
}

func (service depositImpl) Delete(id *string) *model.Erro {
	deposit, err := service.Get(id)
	if err != nil {
		return err
	}
	account, err := service.accountService.Get(&deposit.Account_id)
	if err != nil {
		return err
	}
	account.Balance -= deposit.Deposit

	if _, err := service.accountService.Update(account); err != nil {
		return err
	} else {
		if err := service.depositDatabase.Delete(id); err != nil {
			return err
		}
	}

	return nil
}

func (service depositImpl) Get(id *string) (*deposit.Deposit, *model.Erro) {
	obj, err := service.depositDatabase.Get(id)
	if err != nil {
		return nil, err
	}

	depositResponse, ok := obj.(*deposit.Deposit)
	if !ok {
		return nil, model.DataTypeWrong
	}

	return depositResponse, nil
}

func (service depositImpl) GetAll() (*[]deposit.Deposit, *model.Erro) {
	obj, err := service.depositDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	deposits, ok := obj.(*[]deposit.Deposit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return deposits, nil
}

func (service depositImpl) ProcessDeposit(depositRequest *deposit.Deposit) (*deposit.Deposit, *model.Erro) {
	accountRequest, err := service.accountService.Get(&depositRequest.Account_id)
	if err != nil {
		return nil, err
	}
	if ok, err := verifyDeposit(depositRequest, accountRequest); !ok {
		return nil, err
	}
	accountRequest.Balance = accountRequest.Balance + depositRequest.Deposit

	depositResponse, err := service.Create(depositRequest)
	if err != nil {
		return nil, err
	}

	if _, err := service.accountService.Update(accountRequest); err != nil {
		if err := service.Delete(&depositResponse.Deposit_id); err != nil {
			log.Panic().Msg("Error deleting deposit after update failure: " + err.Err.Error())
		}
		return nil, err
	}
	log.Info().Msg("Deposit created: " + depositResponse.Deposit_id)
	return depositResponse, nil
}

func verifyDeposit(depositRequest *deposit.Deposit, accountResponse *account.Account) (bool, *model.Erro) {
	if accountResponse.Client_id != depositRequest.Client_id {
		return false, &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountResponse.User_id != depositRequest.User_id {
		return false, &model.Erro{Err: errors.New("User ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountResponse.Agency_id != depositRequest.Agency_id {
		return false, &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	return true, nil
}
