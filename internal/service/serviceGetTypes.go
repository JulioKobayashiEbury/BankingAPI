package service

import (
	"errors"

	model "BankingAPI/internal/model"

	"github.com/rs/zerolog/log"
)

func getAccount(accountID string) (*AccountDB, error) {
	docSnapshot, err := model.GetTypeFromDB(&accountID, model.AccountsPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	var accountDB AccountDB
	if err := docSnapshot.DataTo(&accountDB); err != nil {
		return nil, err
	}
	return &accountDB, nil
}

func AccountResponse(accountID string) (*model.AccountResponse, error) {
	accountDB, err := getAccount(accountID)
	if err != nil {
		return nil, err
	}
	return &model.AccountResponse{
		Account_id: accountDB.account_id,
		Client_id:  accountDB.client_id,
		User_id:    accountDB.user_id,
		Agency_id:  accountDB.agency_id,
		Balance:    accountDB.balance,
		Status:     accountDB.status,
	}, nil
}

func GetAccountByFilterAndOrder(listRequest *model.ListRequest) (*[]model.AccountResponse, error) {
	return nil, nil
}

func GetClient(clientID string) (*model.ClientResponse, error) {
	docSnapshot, err := model.GetTypeFromDB(&clientID, model.ClientPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	var clientDB ClientDB
	if err := docSnapshot.DataTo(&clientDB); err != nil {
		return nil, err
	}
	return &model.ClientResponse{
		Client_id:     clientDB.client_id,
		User_id:       clientDB.user_id,
		Name:          clientDB.name,
		Document:      clientDB.document,
		Register_date: clientDB.register_date,
		Status:        clientDB.status,
	}, nil
}

func GetUser(userID string) (*model.UserResponse, error) {
	docSnapshot, err := model.GetTypeFromDB(&userID, model.UsersPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	if (*docSnapshot).Exists() {
		return nil, errors.New("Normal")
	}
	//WORK HERE
	/* var UserDB UserDB
	if err := docSnapshot.DataTo(&UserDB); err != nil {
		return nil, err
	}
	return &model.UserResponse{
		User_id:       UserDB.user_id,
		Name:          UserDB.name,
		Document:      UserDB.document,
		Register_date: UserDB.register_date,
		Status:        UserDB.status,
	}, nil */
	return nil, nil
}
