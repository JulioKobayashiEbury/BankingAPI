package service

import (
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

func GenerateReportByAccount(accountID *string) (*account.AccountReport, *model.Erro) {
	accountInfo, err := Account(*accountID)
	if err != nil {
		return nil, err
	}
	transfers, err := GetAllTransfersByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	deposits, err := GetAllDepositsByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	withdrawals, err := GetAllWithdrawalsByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	automaticDebits, err := GetAllAutoDebitsByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	accountReport := account.AccountReport{
		Account_id:       accountInfo.Account_id,
		Client_id:        accountInfo.Client_id,
		Agency_id:        accountInfo.Agency_id,
		Balance:          accountInfo.Balance,
		Register_date:    accountInfo.Register_date,
		Status:           accountInfo.Status,
		Transfers:        *transfers,
		Deposits:         *deposits,
		Withdrawals:      *withdrawals,
		Automatic_Debits: *automaticDebits,
		Report_Date:      time.Now().Format(timeLayout),
	}
	log.Info().Msg("Report generated for account: " + *accountID)
	return &accountReport, nil
}

func GenerateReportByClient(clientID *string) (*client.ClientReport, *model.Erro) {
	clientInfo, err := Client(*clientID)
	if err != nil {
		return nil, err
	}
	accounts, err := GetAccountsByClientID(clientID)
	if err != nil {
		return nil, err
	}
	return &client.ClientReport{
		Client_id:     clientInfo.Client_id,
		User_id:       clientInfo.User_id,
		Name:          clientInfo.Name,
		Document:      clientInfo.Document,
		Register_date: clientInfo.Register_date,
		Status:        clientInfo.Status,
		Accounts:      (*accounts),
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}

func GenerateReportByUser(userId *string) (*user.UserReport, *model.Erro) {
	userInfo, err := User(*userId)
	if err != nil {
		return nil, err
	}
	clients, err := GetClientsByUserID(userId)
	if err != nil {
		return nil, err
	}
	return &user.UserReport{
		User_id:       userInfo.User_id,
		Name:          userInfo.Name,
		Document:      userInfo.Document,
		Register_date: userInfo.Register_date,
		Status:        userInfo.Status,
		Clients:       *clients,
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}
