package service

import (
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

func Account(accountID string) (*model.AccountResponse, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&accountID, repository.AccountsPath)
	if err != nil {
		return nil, err
	}
	var accountResponse model.AccountResponse
	if err := docSnapshot.DataTo(&accountResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponse.Account_id = accountID

	log.Info().Msg("Account returned: " + accountID)
	return &accountResponse, nil
}

func GetAccountByFilterAndOrder(listRequest *model.ListRequest) (*[]model.AccountResponse, *model.Erro) {
	return nil, nil
}

func Client(clientID string) (*model.ClientResponse, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&clientID, repository.ClientPath)
	if err != nil {
		return nil, err
	}
	var clientResponse model.ClientResponse
	if err := docSnapshot.DataTo(&clientResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	clientResponse.Client_id = clientID

	log.Info().Msg("Client returned: " + clientID)
	return &clientResponse, nil
}

func User(userID string) (*model.UserResponse, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&userID, repository.UsersPath)
	if err != nil {
		return nil, err
	}

	var userResponse model.UserResponse
	if err := docSnapshot.DataTo(&userResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	userResponse.User_id = userID

	log.Info().Msg("User returned: " + userID)
	return &userResponse, nil
}
