package service

import (
	model "BankingAPI/internal/model"
	"BankingAPI/internal/repository"

	"github.com/rs/zerolog/log"
)

func GetAccount(accountID uint32) (*model.AccountResponse, error) {
	docSnapshot, err := repository.GetTypeFromDB(&accountID, "accounts")
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	var accountDB AccountDB
	if err := docSnapshot.DataTo(&accountDB); err != nil {
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

func GetClient(clientID uint32) (*model.ClientResponse, error) {
	docSnapshot, err := repository.GetTypeFromDB(&clientID, "clients")
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

func GetUser(userID uint32) (*model.UserResponse, error) {
	docSnapshot, err := repository.GetTypeFromDB(&userID, "users")
	if err != nil {
		log.Warn().Msg(err.Error())
		return nil, err
	}
	var UserDB UserDB
	if err := docSnapshot.DataTo(&UserDB); err != nil {
		return nil, err
	}
	return &model.UserResponse{
		User_id:       UserDB.user_id,
		Name:          UserDB.name,
		Document:      UserDB.document,
		Register_date: UserDB.register_date,
		Status:        UserDB.status,
	}, nil
}
