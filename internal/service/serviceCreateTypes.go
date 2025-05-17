package service

import (
	model "BankingAPI/internal/model"
	"BankingAPI/internal/repository"
	"time"
)

func CreateAccount(account *model.AccountRequest) (*model.AccountResponse, error) {
	accountDB := AccountDB{
		client_id: account.Client_id,
		user_id:   account.User_id,
		agency_id: account.Agency_id,
		password:  account.Password,
		balance:   0.0,
		status:    true,
	}
	if err := repository.CreateObject(accountDB, repository.AccountsPath, &accountDB.account_id); err != nil {
		return nil, err
	}
	return &model.AccountResponse{
		Account_id:    accountDB.account_id,
		Client_id:     accountDB.client_id,
		User_id:       accountDB.user_id,
		Agency_id:     accountDB.agency_id,
		Register_date: time.Now().String(),
		Balance:       accountDB.balance,
		Status:        accountDB.status,
	}, nil
}

func CreateClient(client *model.ClientRequest) (*model.ClientResponse, error) {
	clientDB := ClientDB{
		user_id:       client.User_id,
		name:          client.Name,
		document:      client.Document,
		password:      client.Password,
		register_date: time.Now().String(),
		status:        true,
	}
	if err := repository.CreateObject(clientDB, repository.ClientPath, &clientDB.client_id); err != nil {
		return nil, err
	}
	return &model.ClientResponse{
		Client_id:     clientDB.client_id,
		User_id:       clientDB.user_id,
		Name:          clientDB.name,
		Document:      clientDB.document,
		Register_date: time.Now().String(),
		Status:        clientDB.status,
	}, nil
}

func CreateUser(user *model.UserRequest) (*model.UserResponse, error) {
	userDB := UserDB{
		name:          user.Name,
		document:      user.Document,
		password:      user.Password,
		register_date: time.Now().String(),
		status:        true,
	}
	if err := repository.CreateObject(userDB, repository.UsersPath, &userDB.user_id); err != nil {
		return nil, err
	}
	return &model.UserResponse{
		User_id:       userDB.user_id,
		Name:          userDB.name,
		Document:      userDB.document,
		Register_date: userDB.register_date,
		Status:        userDB.status,
	}, nil
}
