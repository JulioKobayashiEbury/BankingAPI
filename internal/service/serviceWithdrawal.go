package service

import (
	"errors"
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func ProcessWithdrawal(withdrawalRequest *model.WithdrawalRequest) (*float64, *model.Erro) {
	// monta update
	withdrawal := WithdrawalDB{
		account_id: withdrawalRequest.Account_id,
		client_id:  withdrawalRequest.Client_id,
		user_id:    withdrawalRequest.User_id,
		agency_id:  withdrawalRequest.Agency_iD,
		password:   withdrawalRequest.Password,
		withdrawal: withdrawalRequest.Withdrawal,
		balance:    0.0,
	}
	account, err := Account(withdrawal.account_id)
	if err != nil {
		return nil, err
	}
	if account.Client_id != withdrawal.client_id {
		return nil, &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if account.User_id != withdrawal.user_id {
		return nil, &model.Erro{Err: errors.New("User ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if account.Agency_id != withdrawal.agency_id {
		return nil, &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	/* if account.Password != withdrawal.password {

	}
	*/
	withdrawal.balance = (account.Balance - withdrawal.withdrawal)
	updates := []firestore.Update{
		{
			Path:  "balance",
			Value: withdrawal.balance,
		},
	}

	if err := repository.UpdateTypesDB(&updates, &withdrawal.account_id, repository.AccountsPath); err != nil {
		return nil, err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawal.account_id)
	return &withdrawal.balance, nil
}
