package service

import (
	"errors"
	"fmt"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func UpdateAccount(account *model.AccountRequest) (*model.AccountResponse, error) {
	// verifica valores que foram passados ou não
	paramUpdates := make(map[string]string, 0)
	if account.Account_id == "" {
		log.Warn().Msg("No account with id: 0 allowed")
		return nil, errors.New("Account id invalid")
	}
	if account.Agency_id != 0 {
		paramUpdates["agency_id"] = fmt.Sprintf("%v", account.Agency_id)
	}
	if account.Client_id != "" {
		paramUpdates["client_id"] = fmt.Sprintf("%v", account.Client_id)
	}
	if account.User_id != "" {
		paramUpdates["user_id"] = fmt.Sprintf("%v", account.User_id)
	}
	if account.Password != "" {
		paramUpdates["password"] = account.Password
	}
	// monta struct de update
	updates := make([]firestore.Update, 0, 0)
	for key, value := range paramUpdates {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	if err := model.UpdateTypesDB(&updates, &((*account).Account_id), model.AccountsPath); err != nil {
		return nil, err
	}

	return AccountResponse((*account).Account_id)
}

func UpdateClient(client *model.ClientRequest) (*model.ClientResponse, error) {
	paramUpdates := make(map[string]string, 0)
	if client.Client_id == "" {
		log.Warn().Msg("No client with id: 0 allowed")
		return nil, errors.New("Client id invalid")
	}
	if client.User_id != "" {
		paramUpdates["user_id"] = fmt.Sprintf("%v", client.User_id)
	}
	if client.Name != "" {
		paramUpdates["name"] = client.Name
	}
	if client.Document != "" {
		paramUpdates["document"] = client.Document
	}
	if client.Password != "" {
		paramUpdates["password"] = client.Password
	}
	// monta struct de update
	updates := make([]firestore.Update, 0, 0)
	for key, value := range paramUpdates {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	if err := model.UpdateTypesDB(&updates, &((*client).Client_id), model.ClientPath); err != nil {
		return nil, err
	}

	return GetClient((*client).Client_id)
}

func UpdateUser(user *model.UserRequest) (*model.UserResponse, error) {
	paramUpdates := make(map[string]string, 0)
	if user.User_id == "" {
		log.Warn().Msg("No user with id: 0 allowed")
		return nil, errors.New("User id invalid")
	}
	if user.Name != "" {
		paramUpdates["name"] = user.Name
	}
	if user.Document != "" {
		paramUpdates["document"] = user.Document
	}
	if user.Password != "" {
		paramUpdates["password"] = user.Password
	}
	// monta struct de update
	updates := make([]firestore.Update, 0, 0)
	for key, value := range paramUpdates {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	if err := model.UpdateTypesDB(&updates, &((*user).User_id), model.UsersPath); err != nil {
		return nil, err
	}

	return GetUser((*user).User_id)
}
