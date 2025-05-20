package service

import (
	"errors"
	"fmt"
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func UpdateAccount(account *model.AccountRequest) (*model.AccountResponse, *model.Erro) {
	// verifica valores que foram passados ou não
	paramUpdates := make(map[string]interface{}, 0)
	if account.Account_id == "" {
		log.Warn().Msg("No account with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("Account id invalid"), HttpCode: http.StatusBadRequest}
	}
	if account.Agency_id != 0 {
		paramUpdates["agency_id"] = account.Agency_id
	}
	if account.Client_id != "" {
		paramUpdates["client_id"] = account.Client_id
	}
	if account.User_id != "" {
		paramUpdates["user_id"] = account.User_id
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

	if err := repository.UpdateTypesDB(&updates, &((*account).Account_id), repository.AccountsPath); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (account): " + account.Account_id)
	return Account((*account).Account_id)
}

func UpdateClient(client *model.ClientRequest) (*model.ClientResponse, *model.Erro) {
	paramUpdates := make(map[string]string, 0)
	if client.Client_id == "" {
		log.Warn().Msg("No client with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("Client id invalid"), HttpCode: http.StatusBadRequest}
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

	if err := repository.UpdateTypesDB(&updates, &((*client).Client_id), repository.ClientPath); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (client): " + client.Client_id)
	return Client((*client).Client_id)
}

func UpdateUser(user *model.UserRequest) (*model.UserResponse, *model.Erro) {
	paramUpdates := make(map[string]string, 0)
	if user.User_id == "" {
		log.Warn().Msg("No user with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("User id invalid"), HttpCode: http.StatusBadRequest}
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

	if err := repository.UpdateTypesDB(&updates, &((*user).User_id), repository.UsersPath); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (user): " + user.User_id)

	return User((*user).User_id)
}
