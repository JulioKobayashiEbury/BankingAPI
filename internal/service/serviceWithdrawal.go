package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/withdrawal"

	"github.com/rs/zerolog/log"
)

type withdrawalImpl struct {
	withdrawalDatabase model.RepositoryInterface
	accountService     AccountService
}

func NewWithdrawalService(withdrawalDB model.RepositoryInterface, accountServe AccountService) WithdrawalService {
	return withdrawalImpl{
		withdrawalDatabase: withdrawalDB,
		accountService:     accountServe,
	}
}

func (service withdrawalImpl) Create(withdrawalRequest *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *model.Erro) {
	obj, err := service.withdrawalDatabase.Create(withdrawalRequest)
	if err != nil {
		return nil, err
	}
	withdrawalResponse, ok := obj.(*withdrawal.Withdrawal)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return withdrawalResponse, nil
}

func (service withdrawalImpl) Delete(id *string) *model.Erro {
	// modificar para dar rollback
	withdrawalResponse, err := service.Get(id)
	if err != nil {
		return err
	}

	accountResponse, err := service.accountService.Get(&withdrawalResponse.Account_id)
	if err != nil {
		return err
	}

	accountResponse.Balance += withdrawalResponse.Withdrawal

	if _, err := service.accountService.Update(accountResponse); err != nil {
		return err
	} else {
		if err := service.withdrawalDatabase.Delete(id); err != nil {
			return err
		}
	}

	return nil
}

func (service withdrawalImpl) Get(id *string) (*withdrawal.Withdrawal, *model.Erro) {
	obj, err := service.withdrawalDatabase.Get(id)
	if err != nil {
		return nil, err
	}

	withdrawalResponse, ok := obj.(*withdrawal.Withdrawal)
	if !ok {
		return nil, model.DataTypeWrong
	}

	return withdrawalResponse, nil
}

func (service withdrawalImpl) GetAll() (*[]withdrawal.Withdrawal, *model.Erro) {
	obj, err := service.withdrawalDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	withdrawals, ok := obj.(*[]withdrawal.Withdrawal)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return withdrawals, nil
}

func (service withdrawalImpl) ProcessWithdrawal(withdrawalRequest *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *model.Erro) {
	// monta update
	accountResponse, err := service.accountService.Get(&withdrawalRequest.Account_id)
	if err != nil {
		return nil, err
	}

	if ok, err := verifyWithdrawal(withdrawalRequest, accountResponse); !ok {
		return nil, err
	}

	withdrawalResponse, err := service.Create(withdrawalRequest)
	if err != nil {
		return nil, err
	}

	accountResponse.Balance = accountResponse.Balance - withdrawalRequest.Withdrawal

	if _, err := service.accountService.Update(accountResponse); err != nil {
		if err := service.Delete(&withdrawalResponse.Withdrawal_id); err != nil {
			log.Error().Msg("Account and Withdrawals DB changes failed during processing withdrawal")
			return nil, err
		}
		log.Error().Msg("Creating Account Update failed, withdrawal reversed")
		return nil, err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawalResponse.Account_id)
	return withdrawalResponse, nil
}

func verifyWithdrawal(withdrawalRequest *withdrawal.Withdrawal, accountResponse *account.Account) (bool, *model.Erro) {
	if accountResponse.Client_id != withdrawalRequest.Client_id {
		return false, &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}

	if accountResponse.Agency_id != withdrawalRequest.Agency_id {
		return false, &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountResponse.Balance-withdrawalRequest.Withdrawal < 0.0 {
		return false, &model.Erro{Err: errors.New("Insuficcient funds"), HttpCode: http.StatusBadRequest}
	}
	return true, nil
}
