package service

import (
	"errors"
	"fmt"

	model "BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
)

func ProcessDeposit(depositRequest *model.DepositRequest) (*float64, error) {
	deposit := DepositDB{
		account_id: depositRequest.Account_id,
		client_id:  depositRequest.Client_id,
		user_id:    depositRequest.User_id,
		agency_id:  depositRequest.Agency_iD,
		deposit:    depositRequest.Deposit,
		balance:    0.0,
	}
	account, err := AccountResponse(deposit.account_id)
	if err != nil {
		return nil, err
	}
	if account.Client_id != deposit.client_id {
		return nil, errors.New("Client ID not valid")
	}
	if account.User_id != deposit.user_id {
		return nil, errors.New("User ID not valid")
	}
	if account.Agency_id != deposit.agency_id {
		return nil, errors.New("Agency ID not valid")
	}
	/* if account.Password != deposit.password {

	}
	*/
	deposit.balance = (account.Balance + deposit.deposit)
	updates := []firestore.Update{
		{
			Path:  "Balance",
			Value: fmt.Sprintf("%v", deposit.balance),
		},
	}

	if err := model.UpdateTypesDB(&updates, &deposit.account_id, model.AccountsPath); err != nil {
		return nil, err
	}
	return &deposit.balance, nil
}
