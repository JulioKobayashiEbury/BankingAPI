package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

func CreateAccount(account *model.AccountRequest) (*model.AccountResponse, *model.Erro) {
	if account.Agency_id == 0 || account.User_id == "" || account.Client_id == "" || account.Password == "" {
		log.Warn().Msg("Missing credentials on creating account")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	// verify if client and user exists, PERMISSION MUST BE of user
	if ok, err := verifyUserId(&account.User_id); !ok {
		return nil, err
	}
	accountMap := map[string]interface{}{
		"client_id":     account.Client_id,
		"user_id":       account.User_id,
		"agency_id":     account.Agency_id,
		"password":      account.Password,
		"balance":       0.0,
		"register_date": time.Now().Format(timeLayout),
		"status":        true,
	}
	if err := repository.CreateObject(&accountMap, repository.AccountsPath, &account.Account_id); err != nil {
		return nil, err
	}
	log.Info().Msg("Account created: " + account.Account_id)
	return Account(account.Account_id)
}

func CreateClient(client *model.ClientRequest) (*model.ClientResponse, *model.Erro) {
	if client.User_id == "" || client.Document == "" || client.Password == "" || client.Name == "" {
		log.Warn().Msg("Missing credentials on creating client")
		return nil, &model.Erro{Err: errors.New("Missing credentials for creating client"), HttpCode: http.StatusBadRequest}
	}
	// verify user id exists, PERMISSION MUST BE of user to create
	if ok, err := verifyUserId(&client.User_id); !ok {
		return nil, err
	}
	clientMap := map[string]interface{}{
		"user_id":       client.User_id,
		"name":          client.Name,
		"document":      client.Document,
		"password":      client.Password,
		"register_date": time.Now().Format(timeLayout),
		"status":        true,
	}
	if err := repository.CreateObject(&clientMap, repository.ClientPath, &client.Client_id); err != nil {
		return nil, err
	}
	log.Info().Msg("Client created: " + client.Client_id)
	return Client(client.Client_id)
}

func CreateUser(user *model.UserRequest) (*model.UserResponse, *model.Erro) {
	if user.Name == "" || user.Document == "" || user.Password == "" {
		log.Warn().Msg("Missing credentials on creating user")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	// to create user, permission must be admin
	userMap := map[string]interface{}{
		"name":          user.Name,
		"document":      user.Document,
		"password":      user.Password,
		"register_date": time.Now().Format(timeLayout),
		"status":        true,
	}
	fmt.Println(time.Now().String())
	if err := repository.CreateObject(&userMap, repository.UsersPath, &user.User_id); err != nil {
		return nil, err
	}
	log.Info().Msg("User created: " + user.User_id)
	return User(user.User_id)
}

func verifyUserId(userID *string) (bool, *model.Erro) {
	userSnapshot, err := repository.GetTypeFromDB(userID, repository.UsersPath)
	if err != nil {
		return false, err
	}
	expectedID := userSnapshot.Ref.ID
	fmt.Println(expectedID)
	if expectedID != *userID {
		log.Warn().Msg("No match for user id")
		return false, &model.Erro{Err: errors.New("No match for user id"), HttpCode: http.StatusBadRequest}
	}
	return true, nil
}
