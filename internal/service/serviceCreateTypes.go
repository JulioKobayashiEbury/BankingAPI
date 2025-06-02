package service

import (
	"errors"
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

func CreateAccount(accountRequest *account.AccountRequest) (*account.AccountResponse, *model.Erro) {
	if accountRequest.User_id == "" || accountRequest.Client_id == "" {
		log.Warn().Msg("Missing credentials on creating account")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	// verify if client and user exists, PERMISSION MUST BE of user
	userDatabase := &user.UserFirestore{}
	userDatabase.Request = &user.UserRequest{
		User_id: accountRequest.User_id,
	}
	if err := userDatabase.Get(); err == model.IDnotFound || err != nil {
		return nil, err
	}

	clientDatabase := &client.ClientFirestore{}
	clientDatabase.Request = &client.ClientRequest{
		Client_id: accountRequest.Client_id,
	}
	if err := clientDatabase.Get(); err == model.IDnotFound || err != nil {
		return nil, err
	}

	accountDatabase := &account.AccountFirestore{
		Request: accountRequest,
	}
	if err := accountDatabase.Create(); err != nil {
		return nil, err
	}
	accountRequest.Account_id = accountDatabase.Response.Account_id

	log.Info().Msg("Account created: " + accountRequest.Account_id)
	return Account(accountRequest.Account_id)
}

func CreateClient(clientRequest *client.ClientRequest) (*client.ClientResponse, *model.Erro) {
	if clientRequest.User_id == "" || clientRequest.Document == "" || clientRequest.Name == "" {
		log.Warn().Msg("Missing credentials on creating client")
		return nil, &model.Erro{Err: errors.New("Missing credentials for creating client"), HttpCode: http.StatusBadRequest}
	}

	userDatabase := &user.UserFirestore{}
	userDatabase.Request = &user.UserRequest{
		User_id: clientRequest.User_id,
	}
	if err := userDatabase.Get(); err == model.IDnotFound || err != nil {
		return nil, err
	}
	// verify user id exists, PERMISSION MUST BE of user to create
	log.Info().Msg("Client created: " + clientRequest.Client_id)

	clientDatabase := &client.ClientFirestore{
		Request: clientRequest,
	}
	if err := clientDatabase.Create(); err != nil {
		return nil, err
	}
	clientRequest.Client_id = clientDatabase.Response.Client_id
	return Client(clientRequest.Client_id)
}

func CreateUser(userRequest *user.UserRequest) (*user.UserResponse, *model.Erro) {
	if userRequest.Name == "" || userRequest.Document == "" || userRequest.Password == "" {
		log.Warn().Msg("Missing credentials on creating user")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	userDatabase := &user.UserFirestore{
		Request: userRequest,
	}
	if err := userDatabase.Create(); err != nil {
		return nil, err
	}
	userRequest.User_id = userDatabase.Response.User_id
	log.Info().Msg("User created: " + userRequest.User_id)
	return User(userRequest.User_id)
}
