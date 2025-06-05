package service

import (
	"errors"
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

func (create createImpl) CreateClient(Client *client.Client) (*client.Client, *model.Erro) {
	if Client.User_id == "" || Client.Document == "" || Client.Name == "" {
		log.Warn().Msg("Missing credentials on creating client")
		return nil, &model.Erro{Err: errors.New("Missing credentials for creating client"), HttpCode: http.StatusBadRequest}
	}

	if _, err := create.userDatabase.Get(&Client.User_id); err == model.IDnotFound || err != nil {
		return nil, err
	}
	// verify user id exists, PERMISSION MUST BE of user to create
	log.Info().Msg("Client created: " + Client.Client_id)
	clientID, err := create.clientDatabase.Create(Client)
	if err != nil {
		return nil, err
	}
	return create.getService.Client(*clientID)
}

func (create createImpl) CreateUser(userRequest *user.User) (*user.User, *model.Erro) {
	if userRequest.Name == "" || userRequest.Document == "" || userRequest.Password == "" {
		log.Warn().Msg("Missing credentials on creating user")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	userID, err := create.userDatabase.Create(userRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("User created: " + *userID)
	return create.getService.User(*userID)
}
