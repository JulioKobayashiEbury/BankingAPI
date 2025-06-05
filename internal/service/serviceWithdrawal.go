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
	accountDatabase    model.RepositoryInterface
	withdrawalDatabase model.RepositoryInterface
	getService         ServiceGet
}

func NewWithdrawalService(accountDB model.RepositoryInterface, withdrawalDB model.RepositoryInterface, get ServiceGet) WithdrawalService {
	return withdrawalImpl{
		accountDatabase:    accountDB,
		withdrawalDatabase: withdrawalDB,
		getService:         get,
	}
}

func (wihdrawal withdrawalImpl) Create(withdrawalRequest *withdrawal.Withdrawal) (*string, *model.Erro) {
	return nil, nil
}

func (wihdrawal withdrawalImpl) Delete(*string) *model.Erro {
	return nil
}

func (wihdrawal withdrawalImpl) GetAll(*string) ([]*withdrawal.Withdrawal, *model.Erro) {
	return nil, nil
}

func (withdrawal withdrawalImpl) ProcessWithdrawal(withdrawalRequest *withdrawal.Withdrawal) (*string, *model.Erro) {
	// monta update
	accountResponse, err := withdrawal.getService.Account(withdrawalRequest.Account_id)
	if err != nil {
		return nil, err
	}

	if ok, err := verifyWithdrawal(withdrawalRequest, accountResponse); !ok {
		return nil, err
	}

	withdrawalID, err := withdrawal.withdrawalDatabase.Create(withdrawalRequest)
	if err != nil {
		return nil, err
	}

	accountResponse.Balance = accountResponse.Balance - withdrawalRequest.Withdrawal

	if err := withdrawal.accountDatabase.Update(accountResponse); err != nil {
		if err := withdrawal.withdrawalDatabase.Delete(withdrawalID); err != nil {
			log.Error().Msg("Account and Withdrawals DB changes failed during processing withdrawal")
			return nil, err
		}
		log.Error().Msg("Creating Account Update failed, withdrawal reversed")
		return nil, err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawalRequest.Account_id)
	return withdrawalID, nil
}

func verifyWithdrawal(withdrawalRequest *withdrawal.Withdrawal, accountResponse *account.Account) (bool, *model.Erro) {
	if accountResponse.Client_id != withdrawalRequest.Client_id {
		return false, &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}

	if accountResponse.Agency_id != withdrawalRequest.Agency_id {
		return false, &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	accountResponse.Balance = (accountResponse.Balance - withdrawalRequest.Withdrawal)
	if accountResponse.Balance < 0.0 {
		return false, &model.Erro{Err: errors.New("Insuficcient funds"), HttpCode: http.StatusBadRequest}
	}
	return true, nil
}
