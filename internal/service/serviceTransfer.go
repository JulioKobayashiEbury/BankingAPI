package service

import (
	"fmt"

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
	accountFrom, err := getAccount(transferDBT.account_id_from)
	if err != nil {
		return err
	}

	accountTo, err := getAccount(transferDBT.account_id_to)
	if err != nil {
		return err
	}

	accountTo.balance += transferDBT.value
	accountFrom.balance -= transferDBT.value
	// update from account
	updateOnFrom := []firestore.Update{
		{
			Path:  "Balance",
			Value: fmt.Sprintf("%v", accountFrom.balance),
		},
	}
	updateOnTo := []firestore.Update{
		{
			Path:  "Balance",
			Value: fmt.Sprintf("%v", accountTo.balance),
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
