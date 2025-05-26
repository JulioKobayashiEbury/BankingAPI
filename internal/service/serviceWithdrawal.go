package service

import (
	"errors"
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func ProcessWithdrawal(withdrawalRequest *model.WithdrawalRequest) *model.Erro {
	// monta update
	withdrawal := WithdrawalDB{
		account_id: withdrawalRequest.Account_id,
		client_id:  withdrawalRequest.Client_id,
		agency_id:  withdrawalRequest.Agency_iD,
		withdrawal: withdrawalRequest.Withdrawal,
		balance:    0.0,
	}
	account, err := Account(withdrawal.account_id)
	if err != nil {
		return err
	}
	if account.Client_id != withdrawal.client_id {
		return &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}

	if account.Agency_id != withdrawal.agency_id {
		return &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	withdrawal.balance = (account.Balance - withdrawal.withdrawal)
	if withdrawal.balance < 0.0 {
		return &model.Erro{Err: errors.New("Insuficcient funds"), HttpCode: http.StatusBadRequest}
	}
	updates := []firestore.Update{
		{
			Path:  "balance",
			Value: withdrawal.balance,
		},
	}

	if err := repository.UpdateTypesDB(&updates, &withdrawal.account_id, repository.AccountsPath); err != nil {
		return err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawal.account_id)
	return nil
}
