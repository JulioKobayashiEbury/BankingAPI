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

type ServiceUpdate interface {
	UpdateAccount(accountRequest *account.Account) (*account.Account, *model.Erro)
	UpdateClient(clientRequest *client.ClientRequest) (*client.ClientResponse, *model.Erro)
	UpdateUser(userRequest *user.User) (*user.User, *model.Erro)
}

type updateImpl struct {
	accountDatabase model.RepositoryInterface
	clientDatabase  model.RepositoryInterface
	userDatabase    model.RepositoryInterface
	get             ServiceGet
}

func NewUpdateService(accountDB model.RepositoryInterface, clientDB model.RepositoryInterface, userDB model.RepositoryInterface, getService ServiceGet) ServiceUpdate {
	return updateImpl{
		accountDatabase: accountDB,
		clientDatabase:  clientDB,
		userDatabase:    userDB,
		get:             getService,
	}
}

func (update updateImpl) UpdateAccount(accountRequest *account.Account) (*account.Account, *model.Erro) {
	accountResponse, err := update.get.Account(accountRequest.Account_id)
	if err != nil {
		return nil, err
	}

	// verifica valores que foram passados ou não
	if accountRequest.Account_id == "" {
		log.Warn().Msg("No account with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("Account id invalid"), HttpCode: http.StatusBadRequest}
	}
	if accountRequest.Agency_id != 0 {
		accountResponse.Agency_id = accountRequest.Agency_id
	}
	if accountRequest.Client_id != "" {
		accountResponse.Client_id = accountRequest.Client_id
	}
	if accountRequest.User_id != "" {
		accountResponse.User_id = accountRequest.User_id
	}
	// monta struct de update

	if err := update.accountDatabase.Update(accountResponse); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (account): " + accountRequest.Account_id)

	return update.get.Account(accountRequest.Account_id)
}

func (update updateImpl) UpdateClient(clientRequest *client.ClientRequest) (*client.ClientResponse, *model.Erro) {
	clientResponse, err := update.get.Client(clientRequest.Client_id)
	if err != nil {
		return nil, err
	}

	if clientRequest.User_id != "" {
		clientResponse.User_id = clientRequest.User_id
	}
	if clientRequest.Name != "" {
		clientResponse.Name = clientRequest.Name
	}
	if clientRequest.Document != "" {
		clientResponse.Document = clientRequest.Document
	}
	// monta struct de update
	if err := update.clientDatabase.Update(clientResponse); err != nil {
		return nil, err
	}
	log.Info().Msg("Update was succesful (client): " + clientRequest.Client_id)
	return update.get.Client(clientRequest.Client_id)
}

func (update updateImpl) UpdateUser(userRequest *user.User) (*user.User, *model.Erro) {
	userResponse, err := update.get.User(userRequest.User_id)
	if err != nil {
		return nil, err
	}

	if userRequest.Name != "" {
		userResponse.Name = userRequest.Name
	}
	if userRequest.Document != "" {
		userResponse.Document = userRequest.Document
	}
	if userRequest.Password != "" {
		userResponse.Password = userRequest.Password
	}
	// monta struct de updat

	if err := update.userDatabase.Update(userResponse); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (user): " + userRequest.User_id)

	return update.get.User(userRequest.User_id)
}
