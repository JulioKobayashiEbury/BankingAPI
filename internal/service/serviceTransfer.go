package service

import (
	model "BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
)

func ProcessNewTransfer(transfer *model.TransferRequest) error {
	transferDBT := TransferDB{
		account_id_from: transfer.Account_id_from,
		account_id_to:   transfer.Account_id_to,
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
	if err := model.UpdateTypesDB(&updateOnFrom, &transferDBT.account_id_from, model.AccountsPath); err != nil {
		return err
	}
	if err := model.UpdateTypesDB(&updateOnTo, &transferDBT.account_id_to, model.AccountsPath); err != nil {
		return err
	}
	return nil
}
