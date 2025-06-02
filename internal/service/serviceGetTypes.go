package service

import (
	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/model/withdrawal"

	"github.com/rs/zerolog/log"
)

func Account(accountID string) (*account.AccountResponse, *model.Erro) {
	database := &account.AccountFirestore{}
	database.Request = &account.AccountRequest{
		Account_id: accountID,
	}
	if err := database.Get(); err != nil {
		return nil, err
	}
	log.Info().Msg("Account returned: " + accountID)
	return database.Response, nil
}

func GetAccountByFilterAndOrder() (*[]account.AccountResponse, *model.Erro) {
	return nil, nil
}

func Client(clientID string) (*client.ClientResponse, *model.Erro) {
	database := &client.ClientFirestore{}
	database.Request = &client.ClientRequest{
		Client_id: clientID,
	}
	if err := database.Get(); err != nil {
		return nil, err
	}
	log.Info().Msg("Client returned: " + clientID)
	return database.Response, nil
}

func User(userID string) (*user.UserResponse, *model.Erro) {
	database := &user.UserFirestore{}
	database.Request = &user.UserRequest{
		User_id: userID,
	}

	if err := database.Get(); err != nil {
		return nil, err
	}
	log.Info().Msg("User returned: " + userID)
	return database.Response, nil
}

func GetAllTransfersByAccountID(accountID *string) (*[]transfer.TransferResponse, *model.Erro) {
	database := &transfer.TransferFirestore{}
	if err := database.GetAll(); err != nil {
		return nil, err
	}
	accountTransferSlice := make([]transfer.TransferResponse, 0, len(*database.Slice))
	for _, transfer := range *database.Slice {
		if transfer.Account_id == *accountID {
			accountTransferSlice = append(accountTransferSlice, *transfer)
		}
	}
	return &accountTransferSlice, nil
}

func GetAllAutoDebitsByAccountID(accountID *string) (*[]automaticdebit.AutomaticDebit, *model.Erro) {
	database := &automaticdebit.AutoDebitFirestore{}
	if err := database.GetAll(); err != nil {
		return nil, err
	}
	autoDebitByAccountSlice := make([]automaticdebit.AutomaticDebit, 0, len(*database.Slice))
	for _, autodebit := range *database.Slice {
		if autodebit.Account_id == *accountID {
			autoDebitByAccountSlice = append(autoDebitByAccountSlice, *autodebit)
		}
	}
	return &autoDebitByAccountSlice, nil
}

func GetAllWithdrawalsByAccountID(accountID *string) (*[]withdrawal.WithdrawalResponse, *model.Erro) {
	database := &withdrawal.WithdrawalFirestore{}
	if err := database.GetAll(); err != nil {
		return nil, err
	}
	accountWithdrawalsSlice := make([]withdrawal.WithdrawalResponse, 0, len(*database.Slice))
	for _, withdrawal := range *database.Slice {
		if withdrawal.Account_id == *accountID {
			accountWithdrawalsSlice = append(accountWithdrawalsSlice, *withdrawal)
		}
	}
	return &accountWithdrawalsSlice, nil
}

func GetAllDepositsByAccountID(accountID *string) (*[]deposit.DepositResponse, *model.Erro) {
	database := &deposit.DepositFirestore{}
	if err := database.GetAll(); err != nil {
		return nil, err
	}
	accountDepositsSlice := make([]deposit.DepositResponse, 0, len(*database.Slice))
	for _, deposit := range *database.Slice {
		if deposit.Account_id == *accountID {
			accountDepositsSlice = append(accountDepositsSlice, *deposit)
		}
	}
	return &accountDepositsSlice, nil
}

func GetAccountsByClientID(clientID *string) (*[]account.AccountResponse, *model.Erro) {
	database := &account.AccountFirestore{}
	if err := database.GetAll(); err != nil {
		return nil, err
	}
	clientAccountsSlice := make([]account.AccountResponse, 0, len(*database.Slice))
	for _, accounts := range *database.Slice {
		if accounts.Client_id == *clientID {
			clientAccountsSlice = append(clientAccountsSlice, *accounts)
		}
	}
	return &clientAccountsSlice, nil
}

func GetClientsByUserID(userID *string) (*[]client.ClientResponse, *model.Erro) {
	return nil, nil
}
