package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/withdrawal"

	"github.com/rs/zerolog/log"
)

func ProcessWithdrawal(withdrawalRequest *withdrawal.WithdrawalRequest) *model.Erro {
	// monta update
	accountResponse, err := Account(withdrawalRequest.Account_id)
	if err != nil {
		return err
	}
	if accountResponse.Client_id != withdrawalRequest.Client_id {
		return &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}

	if accountResponse.Agency_id != withdrawalRequest.Agency_id {
		return &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	accountResponse.Balance = (accountResponse.Balance - withdrawalRequest.Withdrawal)
	if accountResponse.Balance < 0.0 {
		return &model.Erro{Err: errors.New("Insuficcient funds"), HttpCode: http.StatusBadRequest}
	}
	databaseWithdrawal := &withdrawal.WithdrawalFirestore{Request: withdrawalRequest}
	if err := databaseWithdrawal.Create(); err != nil {
		return err
	}

	databaseAccount := &account.AccountFirestore{
		Request: &account.AccountRequest{
			Account_id: withdrawalRequest.Account_id,
		},
	}

	databaseAccount.AddUpdate("balance", accountResponse.Balance)

	if err := databaseAccount.Update(); err != nil {
		if err := databaseWithdrawal.Delete(); err != nil {
			log.Error().Msg("Account and Withdrawals DB changes failed during processing withdrawal")
			return err
		}
		log.Error().Msg("Creating Account Update failed, withdrawal reversed")
		return err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawalRequest.Account_id)
	return nil
}
