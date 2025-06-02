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
	database.Request.Account_id = accountID

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
	database.Request.Client_id = clientID

	if err := database.Get(); err != nil {
		return nil, err
	}
	log.Info().Msg("Account returned: " + clientID)
	return database.Response, nil
}

func User(userID string) (*user.UserResponse, *model.Erro) {
	database := &user.UserFirestore{}
	database.Request.User_id = userID

	if err := database.Get(); err != nil {
		return nil, err
	}
	log.Info().Msg("Account returned: " + userID)
	return database.Response, nil
}

func GetAllTransfers(accounID *string) (*[]transfer.TransferResponse, *model.Erro) {
	return nil, nil
}

func GetAllAutoDebits(accountID *string) (*[]automaticdebit.AutomaticDebit, *model.Erro) {
	return nil, nil
}

func GetAllWithdrawals(accountID *string) (*[]withdrawal.WithdrawalResponse, *model.Erro) {
	return nil, nil
}

func GetAllDeposits(accountID *string) (*[]deposit.DepositResponse, *model.Erro) {
	return nil, nil
}

func GetAccountsByClientID(clientID *string) (*[]account.AccountResponse, *model.Erro) {
	return nil, nil
}

func GetClientsByUserID(userID *string) (*[]client.ClientResponse, *model.Erro) {
	return nil, nil
}
