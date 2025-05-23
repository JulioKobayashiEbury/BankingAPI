package service

import (
	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func ProcessNewTransfer(transfer *model.TransferRequest) *model.Erro {
	transferDBT := TransferDB{
		account_id_from: transfer.Account_id,
		account_id_to:   transfer.Account_to,
		value:           transfer.Value,
		password:        transfer.Password,
	}
	// get account from and authenticate
	accountFrom, err := Account(transferDBT.account_id_from)
	if err != nil {
		return err
	}

	accountTo, err := Account(transferDBT.account_id_to)
	if err != nil {
		return err
	}

	accountFrom.Balance -= transferDBT.value
	accountTo.Balance += transferDBT.value

	updateOnFrom := []firestore.Update{
		{
			Path:  "balance",
			Value: accountFrom.Balance,
		},
	}
	updateOnTo := []firestore.Update{
		{
			Path:  "balance",
			Value: accountTo.Balance,
		},
	}
	if err := repository.UpdateTypesDB(&updateOnFrom, &transferDBT.account_id_from, repository.AccountsPath); err != nil {
		return err
	}
	if err := repository.UpdateTypesDB(&updateOnTo, &transferDBT.account_id_to, repository.AccountsPath); err != nil {
		return err
	}
	transferMap := map[string]interface{}{
		"account_id": accountFrom.Account_id,
		"account_to": accountTo.Account_id,
		"value":      transferDBT.value,
	}
	var transferID string
	if err := repository.CreateObject(&transferMap, repository.TransfersPath, &transferID); err != nil {
		return err
	}
	(*transfer).Transfer_id = transferID

	log.Info().Msg("Transfer was succesful: " + transferDBT.account_id_from + " to " + transferDBT.account_id_to)
	return nil
}

func rollBackTranfer() {
}
