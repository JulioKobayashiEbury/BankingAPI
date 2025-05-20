package service

import (
	"time"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

func CreateAccount(account *model.AccountRequest) (*model.AccountResponse, *model.Erro) {
	accountMap := map[string]interface{}{
		"client_id":     account.Client_id,
		"user_id":       account.User_id,
		"agency_id":     account.Agency_id,
		"password":      account.Password,
		"balance":       0.0,
		"register_date": time.Now().String(),
		"status":        true,
	}
	if err := repository.CreateObject(&accountMap, repository.AccountsPath, &account.Account_id); err != nil {
		return nil, err
	}
	log.Info().Msg("Account created: " + account.Account_id)
	return Account(account.Account_id)
}

func CreateClient(client *model.ClientRequest) (*model.ClientResponse, *model.Erro) {
	clientMap := map[string]interface{}{
		"user_id":       client.User_id,
		"name":          client.Name,
		"document":      client.Document,
		"password":      client.Password,
		"register_date": time.Now().String(),
		"status":        true,
	}
	if err := repository.CreateObject(&clientMap, repository.ClientPath, &client.Client_id); err != nil {
		return nil, err
	}
	log.Info().Msg("Client created: " + client.Client_id)
	return Client(client.Client_id)
}

func CreateUser(user *model.UserRequest) (*model.UserResponse, *model.Erro) {
	userMap := map[string]interface{}{
		"name":          user.Name,
		"document":      user.Document,
		"password":      user.Password,
		"register_date": time.Now().String(),
		"status":        true,
	}
	if err := repository.CreateObject(&userMap, repository.UsersPath, &user.User_id); err != nil {
		return nil, err
	}
	log.Info().Msg("User created: " + user.User_id)
	return User(user.User_id)
}
