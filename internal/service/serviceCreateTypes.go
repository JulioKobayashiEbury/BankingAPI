package service

import (
	"fmt"
	"time"

	model "BankingAPI/internal/model"
)

func CreateAccount(account *model.AccountRequest) (*model.AccountResponse, error) {
	accountMap := map[string]interface{}{
		"client_id":     account.Client_id,
		"user_id":       account.User_id,
		"agency_id":     account.Agency_id,
		"password":      account.Password,
		"balance":       0.0,
		"register_date": time.Now().String(),
		"status":        true,
	}
	if err := model.CreateObject(&accountMap, model.AccountsPath, &account.Account_id); err != nil {
		return nil, err
	}
	return Account(account.Account_id)
}

func CreateClient(client *model.ClientRequest) (*model.ClientResponse, error) {
	clientMap := map[string]interface{}{
		"user_id":       client.User_id,
		"name":          client.Name,
		"document":      client.Document,
		"password":      client.Password,
		"register_date": time.Now().String(),
		"status":        true,
	}
	if err := model.CreateObject(&clientMap, model.ClientPath, &client.Client_id); err != nil {
		return nil, err
	}
	return Client(client.Client_id)
}

func CreateUser(user *model.UserRequest) (*model.UserResponse, error) {
	userMap := map[string]interface{}{
		"name":          user.Name,
		"document":      user.Document,
		"password":      user.Password,
		"register_date": time.Now().String(),
		"status":        true,
	}
	if err := model.CreateObject(&userMap, model.UsersPath, &user.User_id); err != nil {
		return nil, err
	}
	fmt.Print(userMap["id"])
	return User(user.User_id)
}
