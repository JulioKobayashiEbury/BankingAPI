package service

import (
	"errors"
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func ProcessDeposit(depositRequest *model.DepositRequest) (*float64, *model.Erro) {
	deposit := DepositDB{
		account_id: (*depositRequest).Account_id,
		client_id:  (*depositRequest).Client_id,
		user_id:    (*depositRequest).User_id,
		agency_id:  (*depositRequest).Agency_iD,
		deposit:    (*depositRequest).Deposit,
		balance:    0.0,
	}
	account, err := Account(deposit.account_id)
	if err != nil {
		return nil, err
	}
	if account.Client_id != deposit.client_id {
		return nil, &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if account.User_id != deposit.user_id {
		return nil, &model.Erro{Err: errors.New("User ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if account.Agency_id != deposit.agency_id {
		return nil, &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}

	deposit.balance = ((*account).Balance + deposit.deposit)
	updates := []firestore.Update{
		{
			Path:  "balance",
			Value: deposit.balance,
		},
	}
	if err := repository.UpdateTypesDB(&updates, &deposit.account_id, repository.AccountsPath); err != nil {
		return nil, err
	}
	log.Info().Msg("Deposit created: " + deposit.account_id)
	return &deposit.balance, nil
}
