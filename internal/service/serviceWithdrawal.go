package service

import (
	"context"
	"errors"
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/withdrawal"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type withdrawalImpl struct {
	withdrawalDatabase withdrawal.WithdrawalRepository
	accountService     AccountService
}

func NewWithdrawalService(withdrawalDB withdrawal.WithdrawalRepository, accountServe AccountService) WithdrawalService {
	return withdrawalImpl{
		withdrawalDatabase: withdrawalDB,
		accountService:     accountServe,
	}
}

func (service withdrawalImpl) Create(ctx context.Context, withdrawalRequest *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *echo.HTTPError) {
	withdrawalResponse, err := service.withdrawalDatabase.Create(ctx, withdrawalRequest)
	if err != nil {
		return nil, err
	}
	return withdrawalResponse, nil
}

func (service withdrawalImpl) Delete(ctx context.Context, id *string) *echo.HTTPError {
	// modificar para dar rollback
	withdrawalResponse, err := service.Get(ctx, id)
	if err != nil {
		return err
	}

	accountResponse, err := service.accountService.Get(ctx, &withdrawalResponse.Account_id)
	if err != nil {
		return err
	}

	accountResponse.Balance += withdrawalResponse.Withdrawal

	if _, err := service.accountService.Update(ctx, accountResponse); err != nil {
		return err
	} else {
		if err := service.withdrawalDatabase.Delete(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (service withdrawalImpl) Get(ctx context.Context, id *string) (*withdrawal.Withdrawal, *echo.HTTPError) {
	withdrawalResponse, err := service.withdrawalDatabase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return withdrawalResponse, nil
}

func (service withdrawalImpl) GetAll(ctx context.Context) (*[]withdrawal.Withdrawal, *echo.HTTPError) {
	withdrawals, err := service.withdrawalDatabase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (service withdrawalImpl) ProcessWithdrawal(ctx context.Context, withdrawalRequest *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *echo.HTTPError) {
	// monta update
	accountResponse, err := service.accountService.Get(ctx, &withdrawalRequest.Account_id)
	if err != nil {
		return nil, err
	}

	if ok, err := verifyWithdrawal(withdrawalRequest, accountResponse); !ok {
		return nil, err
	}

	withdrawalResponse, err := service.Create(ctx, withdrawalRequest)
	if err != nil {
		return nil, err
	}

	accountResponse.Balance = accountResponse.Balance - withdrawalRequest.Withdrawal

	if _, err := service.accountService.Update(ctx, accountResponse); err != nil {
		if err := service.Delete(ctx, &withdrawalResponse.Withdrawal_id); err != nil {
			log.Error().Msg("Account and Withdrawals DB changes failed during processing withdrawal")
			return nil, err
		}
		log.Error().Msg("Creating Account Update failed, withdrawal reversed")
		return nil, err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawalResponse.Account_id)
	return withdrawalResponse, nil
}

func verifyWithdrawal(withdrawalRequest *withdrawal.Withdrawal, accountResponse *account.Account) (bool, *echo.HTTPError) {
	if accountResponse.Status != "active" {
		return false, model.ErrAccountNotActive
	}
	if accountResponse.User_id != withdrawalRequest.User_id {
		return false, model.ErrUserIDNotValid
	}

	if accountResponse.Agency_id != withdrawalRequest.Agency_id {
		return false, model.ErrAgencyIDNotValid
	}
	if accountResponse.Balance-withdrawalRequest.Withdrawal < 0.0 {
		return false, &echo.HTTPError{Internal: errors.New("insuficcient funds"), Code: http.StatusBadRequest, Message: "insuficcient funds"}
	}
	return true, nil
}
