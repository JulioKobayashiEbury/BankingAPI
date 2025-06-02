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

func UpdateAccount(accountRequest *account.AccountRequest) (*account.AccountResponse, *model.Erro) {
	// verifica valores que foram passados ou não
	database := &account.AccountFirestore{}
	if accountRequest.Account_id == "" {
		log.Warn().Msg("No account with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("Account id invalid"), HttpCode: http.StatusBadRequest}
	}
	if accountRequest.Agency_id != 0 {
		database.AddUpdate("agency_id", accountRequest.Agency_id)
	}
	if accountRequest.Client_id != "" {
		database.AddUpdate("client_id", accountRequest.Client_id)
	}
	if accountRequest.User_id != "" {
		database.AddUpdate("user_id", accountRequest.User_id)
	}
	// monta struct de update

	if err := database.Update(); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (account): " + accountRequest.Account_id)
	return Account(accountRequest.Account_id)
}

func UpdateClient(clientRequest *client.ClientRequest) (*client.ClientResponse, *model.Erro) {
	database := &client.ClientFirestore{}
	database.Request.Client_id = clientRequest.Client_id
	if clientRequest.User_id != "" {
		database.AddUpdate("user_id", clientRequest.User_id)
	}
	if clientRequest.Name != "" {
		database.AddUpdate("name", clientRequest.Name)
	}
	if clientRequest.Document != "" {
		database.AddUpdate("document", clientRequest.Document)
	}
	// monta struct de update
	if err := database.Update(); err != nil {
		return nil, err
	}
	log.Info().Msg("Update was succesful (client): " + clientRequest.Client_id)
	return Client(clientRequest.Client_id)
}

func UpdateUser(userRequest *user.UserRequest) (*user.UserResponse, *model.Erro) {
	database := &user.UserFirestore{}
	database.Request.User_id = userRequest.User_id
	if userRequest.Name != "" {
		database.AddUpdate("name", userRequest.Name)
	}
	if userRequest.Document != "" {
		database.AddUpdate("document", userRequest.Document)
	}
	if userRequest.Password != "" {
		database.AddUpdate("password", userRequest.Password)
	}
	// monta struct de updat

	if err := database.Update(); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (user): " + userRequest.User_id)

	return User(userRequest.User_id)
}
