package service

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/withdrawal"
)

type getFilteredImpl struct {
	clientDatabase     model.RepositoryInterface
	accountDatabse     model.RepositoryInterface
	transferDatabase   model.RepositoryInterface
	depositDatabase    model.RepositoryInterface
	withdrawalDatabase model.RepositoryInterface
	autodebitDatabase  model.RepositoryInterface
}

func NewGetFilteredService(
	clientDB model.RepositoryInterface,
	accountDB model.RepositoryInterface,
	transferDB model.RepositoryInterface,
	depositDB model.RepositoryInterface,
	withdrawalDB model.RepositoryInterface,
	autodebitDB model.RepositoryInterface,
) GetFilteredService {
	return getFilteredImpl{
		clientDatabase:     clientDB,
		accountDatabse:     accountDB,
		transferDatabase:   transferDB,
		depositDatabase:    depositDB,
		withdrawalDatabase: withdrawalDB,
		autodebitDatabase:  autodebitDB,
	}
}

func (service getFilteredImpl) GetAllTransfersByAccountID(accountID *string) (*[]transfer.Transfer, *model.Erro) {
	obj, err := service.transferDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	transfers, ok := obj.(*[]transfer.Transfer)
	if !ok {
		return nil, model.DataTypeWrong
	}

	transfersByAccountSlice := make([]transfer.Transfer, 0, len(*transfers))
	for _, transfer := range *transfers {
		if transfer.Account_id == *accountID {
			transfersByAccountSlice = append(transfersByAccountSlice, transfer)
		}
	}
	return &transfersByAccountSlice, nil
}

func (service getFilteredImpl) GetAllAutoDebitsByAccountID(accountID *string) (*[]automaticdebit.AutomaticDebit, *model.Erro) {
	obj, err := service.autodebitDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	autodebits, ok := obj.(*[]automaticdebit.AutomaticDebit)
	if !ok {
		return nil, model.DataTypeWrong
	}

	autodebitByAccountSlice := make([]automaticdebit.AutomaticDebit, 0, len(*autodebits))
	for _, autodebit := range *autodebits {
		if autodebit.Account_id == *accountID {
			autodebitByAccountSlice = append(autodebitByAccountSlice, autodebit)
		}
	}
	return &autodebitByAccountSlice, nil
}

func (service getFilteredImpl) GetAllDepositsByAccountID(accountID *string) (*[]deposit.Deposit, *model.Erro) {
	obj, err := service.depositDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	deposits, ok := obj.(*[]deposit.Deposit)
	if !ok {
		return nil, model.DataTypeWrong
	}
	depositByAccountSlice := make([]deposit.Deposit, 0, len(*deposits))
	for _, deposit := range *deposits {
		if deposit.Account_id == *accountID {
			depositByAccountSlice = append(depositByAccountSlice, deposit)
		}
	}
	return &depositByAccountSlice, nil
}

func (service getFilteredImpl) GetAllWithdrawalsByAccountID(accountID *string) (*[]withdrawal.Withdrawal, *model.Erro) {
	obj, err := service.withdrawalDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	withdrawals, ok := obj.(*[]withdrawal.Withdrawal)
	if !ok {
		return nil, model.DataTypeWrong
	}

	withdrawalByAccountSlice := make([]withdrawal.Withdrawal, 0, len(*withdrawals))
	for _, withdrawal := range *withdrawals {
		if withdrawal.Account_id == *accountID {
			withdrawalByAccountSlice = append(withdrawalByAccountSlice, withdrawal)
		}
	}
	return &withdrawalByAccountSlice, nil
}

func (service getFilteredImpl) GetClientsByUserID(userID *string) (*[]client.Client, *model.Erro) {
	obj, err := service.clientDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	clients, ok := obj.(*[]client.Client)
	if !ok {
		return nil, model.DataTypeWrong
	}

	clientsByUserIDSlice := make([]client.Client, 0, len(*clients))
	for _, client := range *clients {
		if client.User_id == *userID {
			clientsByUserIDSlice = append(clientsByUserIDSlice, client)
		}
	}
	return &clientsByUserIDSlice, nil
}

func (service getFilteredImpl) GetAccountsByClientID(clientID *string) (*[]account.Account, *model.Erro) {
	obj, err := service.accountDatabse.GetAll()
	if err != nil {
		return nil, err
	}
	accountSlice, ok := obj.(*[]account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	clientAccountsSlice := make([]account.Account, 0, len(*accountSlice))
	for _, accounts := range *accountSlice {
		if accounts.Client_id == *clientID {
			clientAccountsSlice = append(clientAccountsSlice, accounts)
		}
	}
	return &clientAccountsSlice, nil
}
