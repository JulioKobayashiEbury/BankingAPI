package service

import (
	model "BankingAPI/internal/model"

	"github.com/rs/zerolog/log"
)

func Account(accountID string) (*model.AccountResponse, error) {
	docSnapshot, err := model.GetTypeFromDB(&accountID, model.AccountsPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	var accountResponse model.AccountResponse
	if err := docSnapshot.DataTo(&accountResponse); err != nil {
		return nil, err
	}
	accountResponse.Account_id = accountID
	return &accountResponse, nil
}

func GetAccountByFilterAndOrder(listRequest *model.ListRequest) (*[]model.AccountResponse, error) {
	return nil, nil
}

func Client(clientID string) (*model.ClientResponse, error) {
	docSnapshot, err := model.GetTypeFromDB(&clientID, model.ClientPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	var clientResponse model.ClientResponse
	if err := docSnapshot.DataTo(&clientResponse); err != nil {
		return nil, err
	}
	clientResponse.Client_id = clientID
	return &clientResponse, nil
}

func User(userID string) (*model.UserResponse, error) {
	docSnapshot, err := model.GetTypeFromDB(&userID, model.UsersPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}

	var userResponse model.UserResponse
	if err := docSnapshot.DataTo(&userResponse); err != nil {
		return nil, err
	}

	userResponse.User_id = userID

	return &userResponse, nil
}
