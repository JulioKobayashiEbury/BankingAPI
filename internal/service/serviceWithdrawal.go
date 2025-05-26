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

func ProcessWithdrawal(withdrawalRequest *model.WithdrawalRequest) *model.Erro {
	// monta update
	withdrawal := WithdrawalDB{
		account_id:      withdrawalRequest.Account_id,
		client_id:       withdrawalRequest.Client_id,
		agency_id:       withdrawalRequest.Agency_iD,
		withdrawal:      withdrawalRequest.Withdrawal,
		withdrawal_date: time.Now().Format(timeLayout),
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
	account.Balance = (account.Balance - withdrawal.withdrawal)
	if account.Balance < 0.0 {
		return &model.Erro{Err: errors.New("Insuficcient funds"), HttpCode: http.StatusBadRequest}
	}
	updates := []firestore.Update{
		{
			Path:  "balance",
			Value: account.Balance,
		},
	}
	withdrawalMap := map[string]interface{}{
		"account_id":      withdrawal.account_id,
		"client_id":       withdrawal.client_id,
		"user_id":         withdrawal.user_id,
		"agency_id":       withdrawal.agency_id,
		"withdrawal":      withdrawal.withdrawal,
		"withdrawal_date": withdrawal.withdrawal_date,
	}
	if err := repository.CreateObject(&withdrawalMap, repository.WithdrawalsPath, &withdrawal.withdrawal_id); err != nil {
		return err
	}

	if err := repository.UpdateTypesDB(&updates, &withdrawal.account_id, repository.AccountsPath); err != nil {
		if err := repository.DeleteObject(&withdrawal.withdrawal_id, repository.WithdrawalsPath); err != nil {
			log.Panic().Msg("Error deleting withdrawal after update failure: " + err.Err.Error())
		}
		return err
	}

	log.Info().Msg("Succesful Withdrawal: " + withdrawal.account_id)
	return nil
}
