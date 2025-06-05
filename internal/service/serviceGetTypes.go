package service

import (
	"errors"
	"net/http"

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

var ErrRepositoryNotSet = &model.Erro{Err: errors.New("repository needed not set"), HttpCode: http.StatusInternalServerError}

type ServiceGet interface {
	Account(accountID string) (*account.Account, *model.Erro)
	Client(clientID string) (*client.ClientResponse, *model.Erro)
	User(userID string) (*user.User, *model.Erro)
}

type ServiceGetAll interface {
	GetAllTransfersByAccountID(accountID *string) (*[]transfer.TransferResponse, *model.Erro)
	GetAllAutoDebitsByAccountID(accountID *string) (*[]automaticdebit.AutomaticDebitResponse, *model.Erro)
	GetAllWithdrawalsByAccountID(accountID *string) (*[]withdrawal.WithdrawalResponse, *model.Erro)
	GetAllDepositsByAccountID(accountID *string) (*[]deposit.DepositResponse, *model.Erro)
	GetAccountsByClientID(clientID *string) (*[]account.Account, *model.Erro)
	GetClientsByUserID(userID *string) (*[]client.ClientResponse, *model.Erro)
}

type getImpl struct {
	userDatabase   model.RepositoryInterface
	accontDatabase model.RepositoryInterface
	clientDatabase model.RepositoryInterface
}

func NewGetService(accountDB model.RepositoryInterface, clientDB model.RepositoryInterface, userDB model.RepositoryInterface) ServiceGet {
	return getImpl{
		userDatabase:   userDB,
		clientDatabase: clientDB,
		accontDatabase: accountDB,
	}
}

type getAllImpl struct {
	autodebitDatabase  model.RepositoryInterface
	withdrawalDatabase model.RepositoryInterface
	depositDatabase    model.RepositoryInterface
	transferDatabase   model.RepositoryInterface
	accountDatabase    model.RepositoryInterface
	clientDatabase     model.RepositoryInterface
}

func NewGetAllService(autodebitDB model.RepositoryInterface, withdrawalDB model.RepositoryInterface, depositDB model.RepositoryInterface, transferDB model.RepositoryInterface, accountDB model.RepositoryInterface, clientDB model.RepositoryInterface) ServiceGetAll {
	return getAllImpl{
		autodebitDatabase:  autodebitDB,
		withdrawalDatabase: withdrawalDB,
		depositDatabase:    depositDB,
		transferDatabase:   transferDB,
		accountDatabase:    accountDB,
		clientDatabase:     clientDB,
	}
}

func (get getImpl) Account(accountID string) (*account.Account, *model.Erro) {
	if get.accontDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.accontDatabase.Get(&accountID)
	if err != nil {
		return nil, err
	}
	accountResponse, ok := obj.(*account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Account returned: " + accountID)
	return accountResponse, nil
}

func (get getImpl) Client(clientID string) (*client.ClientResponse, *model.Erro) {
	if get.clientDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.clientDatabase.Get(&clientID)
	if err != nil {
		return nil, err
	}
	clientResponse, ok := obj.(*client.ClientResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Client returned: " + clientID)
	return clientResponse, nil
}

func (get getImpl) User(userID string) (*user.User, *model.Erro) {
	if get.userDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.userDatabase.Get(&userID)
	if err != nil {
		return nil, err
	}
	userResponse, ok := obj.(*user.User)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("User returned: " + userID)
	return userResponse, nil
}

func (get getAllImpl) GetAccountByFilterAndOrder() (*[]account.Account, *model.Erro) {
	return nil, nil
}

func (get getAllImpl) GetAllTransfersByAccountID(accountID *string) (*[]transfer.TransferResponse, *model.Erro) {
	if get.transferDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.transferDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	transferSlice, ok := obj.(*[]*transfer.TransferResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	accountTransferSlice := make([]transfer.TransferResponse, 0, len(*transferSlice))
	for _, transfer := range *transferSlice {
		if transfer.Account_id == *accountID {
			accountTransferSlice = append(accountTransferSlice, *transfer)
		}
	}
	return &accountTransferSlice, nil
}

func (get getAllImpl) GetAllAutoDebitsByAccountID(accountID *string) (*[]automaticdebit.AutomaticDebitResponse, *model.Erro) {
	if get.autodebitDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.autodebitDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	autodebitSlice, ok := obj.(*[]*automaticdebit.AutomaticDebitResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	autoDebitByAccountSlice := make([]automaticdebit.AutomaticDebitResponse, 0, len(*autodebitSlice))
	for _, autodebit := range *autodebitSlice {
		if autodebit.Account_id == *accountID {
			autoDebitByAccountSlice = append(autoDebitByAccountSlice, *autodebit)
		}
	}
	return &autoDebitByAccountSlice, nil
}

func (get getAllImpl) GetAllWithdrawalsByAccountID(accountID *string) (*[]withdrawal.WithdrawalResponse, *model.Erro) {
	if get.withdrawalDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.withdrawalDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	withdrawalSlice, ok := obj.(*[]*withdrawal.WithdrawalResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	accountWithdrawalsSlice := make([]withdrawal.WithdrawalResponse, 0, len(*withdrawalSlice))
	for _, withdrawal := range *withdrawalSlice {
		if withdrawal.Account_id == *accountID {
			accountWithdrawalsSlice = append(accountWithdrawalsSlice, *withdrawal)
		}
	}
	return &accountWithdrawalsSlice, nil
}

func (get getAllImpl) GetAllDepositsByAccountID(accountID *string) (*[]deposit.DepositResponse, *model.Erro) {
	if get.depositDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.depositDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	depositSlice, ok := obj.(*[]*deposit.DepositResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	accountDepositsSlice := make([]deposit.DepositResponse, 0, len(*depositSlice))
	for _, deposit := range *depositSlice {
		if deposit.Account_id == *accountID {
			accountDepositsSlice = append(accountDepositsSlice, *deposit)
		}
	}
	return &accountDepositsSlice, nil
}

func (get getAllImpl) GetAccountsByClientID(clientID *string) (*[]account.Account, *model.Erro) {
	if get.accountDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.accountDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	accountSlice, ok := obj.(*[]*account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	clientAccountsSlice := make([]account.Account, 0, len(*accountSlice))
	for _, accounts := range *accountSlice {
		if accounts.Client_id == *clientID {
			clientAccountsSlice = append(clientAccountsSlice, *accounts)
		}
	}
	return &clientAccountsSlice, nil
}

func (get getAllImpl) GetClientsByUserID(userID *string) (*[]client.ClientResponse, *model.Erro) {
	if get.clientDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.clientDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	clientSlice, ok := obj.(*[]*client.ClientResponse)
	if !ok {
		return nil, model.DataTypeWrong
	}
	userClientsSlice := make([]client.ClientResponse, 0, len(*clientSlice))
	for _, clients := range *clientSlice {
		if clients.User_id == *userID {
			userClientsSlice = append(userClientsSlice, *clients)
		}
	}
	return &userClientsSlice, nil
}
