package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

var ErrRepositoryNotSet = &model.Erro{Err: errors.New("repository needed not set"), HttpCode: http.StatusInternalServerError}

func (get getImpl) Client(clientID string) (*client.Client, *model.Erro) {
	if get.clientDatabase == nil {
		return nil, ErrRepositoryNotSet
	}
	obj, err := get.clientDatabase.Get(&clientID)
	if err != nil {
		return nil, err
	}
	Client, ok := obj.(*client.Client)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Client returned: " + clientID)
	return Client, nil
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
