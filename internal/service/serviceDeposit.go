package service

import (
	"errors"
	"net/http"
	"time"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func ProcessDeposit(depositRequest *model.DepositRequest) *model.Erro {
	deposit := DepositDB{
		account_id:   (*depositRequest).Account_id,
		client_id:    (*depositRequest).Client_id,
		user_id:      (*depositRequest).User_id,
		agency_id:    (*depositRequest).Agency_id,
		deposit:      (*depositRequest).Deposit,
		deposit_date: time.Now().Format(timeLayout),
	}
	account, err := Account(deposit.account_id)
	if err != nil {
		return err
	}
	if account.Client_id != deposit.client_id {
		return &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if account.User_id != deposit.user_id {
		return &model.Erro{Err: errors.New("User ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if account.Agency_id != deposit.agency_id {
		return &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}

	account.Balance = account.Balance + deposit.deposit
	updates := []firestore.Update{
		{
			Path:  "balance",
			Value: account.Balance,
		},
	}
	depositMap := map[string]interface{}{
		"account_id":   deposit.account_id,
		"client_id":    deposit.client_id,
		"user_id":      deposit.user_id,
		"agency_id":    deposit.agency_id,
		"deposit":      deposit.deposit,
		"deposit_date": deposit.deposit_date,
	}
	if err := repository.CreateObject(&depositMap, repository.DepositPath, &deposit.deposit_id); err != nil {
		return err
	}
	if err := repository.UpdateTypesDB(&updates, &deposit.account_id, repository.AccountsPath); err != nil {
		if err := repository.DeleteObject(&deposit.deposit_id, repository.DepositPath); err != nil {
			log.Panic().Msg("Error deleting deposit after update failure: " + err.Err.Error())
		}
		return err
	}
	log.Info().Msg("Deposit created: " + deposit.account_id)
	return nil
}
