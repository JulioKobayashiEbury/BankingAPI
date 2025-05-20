package service

import (
	"errors"

	model "BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
)

func ProcessWithdrawal(withdrawalRequest *model.WithdrawalRequest) (*float64, error) {
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
		return nil, errors.New("Client ID not valid")
	}
	if account.User_id != withdrawal.user_id {
		return nil, errors.New("User ID not valid")
	}
	if account.Agency_id != withdrawal.agency_id {
		return nil, errors.New("Agency ID not valid")
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

	if err := model.UpdateTypesDB(&updates, &withdrawal.account_id, model.AccountsPath); err != nil {
		return nil, err
	}
	return &withdrawal.balance, nil
}
