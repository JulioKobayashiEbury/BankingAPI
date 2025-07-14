package service

import (
	"context"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/deposit"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type depositImpl struct {
	depositDatabase deposit.DepositRepository
	accountService  AccountService
}

func NewDepositService(depositDB deposit.DepositRepository, accountServe AccountService) DepositService {
	return depositImpl{
		depositDatabase: depositDB,
		accountService:  accountServe,
	}
}

func (service depositImpl) Create(ctx context.Context, depositRequest *deposit.Deposit) (*deposit.Deposit, *echo.HTTPError) {
	depositResponse, err := service.depositDatabase.Create(ctx, depositRequest)
	if err != nil {
		return nil, err
	}
	return depositResponse, nil
}

func (service depositImpl) Delete(ctx context.Context, id *string) *echo.HTTPError {
	deposit, err := service.Get(ctx, id)
	if err != nil {
		return err
	}
	account, err := service.accountService.Get(ctx, &deposit.Account_id)
	if err != nil {
		return err
	}
	account.Balance -= deposit.Deposit

	if _, err := service.accountService.Update(ctx, account); err != nil {
		return err
	} else {
		if err := service.depositDatabase.Delete(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (service depositImpl) Get(ctx context.Context, id *string) (*deposit.Deposit, *echo.HTTPError) {
	depositResponse, err := service.depositDatabase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return depositResponse, nil
}

func (service depositImpl) GetAll(ctx context.Context) (*[]deposit.Deposit, *echo.HTTPError) {
	deposits, err := service.depositDatabase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return deposits, nil
}

func (service depositImpl) ProcessDeposit(ctx context.Context, depositRequest *deposit.Deposit) (*deposit.Deposit, *echo.HTTPError) {
	accountRequest, err := service.accountService.Get(ctx, &depositRequest.Account_id)
	if err != nil {
		return nil, err
	}
	if ok, err := verifyDeposit(depositRequest, accountRequest); !ok {
		return nil, err
	}
	accountRequest.Balance = accountRequest.Balance + depositRequest.Deposit

	depositResponse, err := service.Create(ctx, depositRequest)
	if err != nil {
		return nil, err
	}

	if _, err := service.accountService.Update(ctx, accountRequest); err != nil {
		if err := service.Delete(ctx, &depositResponse.Deposit_id); err != nil {
			log.Panic().Msg("Error deleting deposit after update failure: " + err.Error())
		}
		return nil, err
	}
	log.Info().Msg("Deposit created: " + depositResponse.Deposit_id)
	return depositResponse, nil
}

func verifyDeposit(depositRequest *deposit.Deposit, accountResponse *account.Account) (bool, *echo.HTTPError) {
	if accountResponse.Status != "active" {
		return false, model.ErrAccountNotActive
	}
	if accountResponse.Client_id != depositRequest.Client_id {
		return false, model.ErrClientIDNotValid
	}
	if accountResponse.User_id != depositRequest.User_id {
		return false, model.ErrUserIDNotValid
	}
	if accountResponse.Agency_id != depositRequest.Agency_id {
		return false, model.ErrAgencyIDNotValid
	}
	return true, nil
}
