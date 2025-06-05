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

type ServiceCreate interface {
	CreateAccount(accountRequest *account.Account) (*account.Account, *model.Erro)
	CreateClient(clientRequest *client.ClientRequest) (*client.ClientResponse, *model.Erro)
	CreateUser(userRequest *user.User) (*user.User, *model.Erro)
}

type createImpl struct {
	accountDatabase model.RepositoryInterface
	clientDatabase  model.RepositoryInterface
	userDatabase    model.RepositoryInterface
	getService      ServiceGet
}

func NewCreateService(accoutnDB model.RepositoryInterface, clientDB model.RepositoryInterface, userDB model.RepositoryInterface, get ServiceGet) ServiceCreate {
	return createImpl{
		accountDatabase: accoutnDB,
		clientDatabase:  clientDB,
		userDatabase:    userDB,
		getService:      get,
	}
}

func (create createImpl) CreateAccount(accountRequest *account.Account) (*account.Account, *model.Erro) {
	if accountRequest.User_id == "" || accountRequest.Client_id == "" {
		log.Warn().Msg("Missing credentials on creating account")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	// verify if client and user exists, PERMISSION MUST BE of user
	if _, err := create.userDatabase.Get(&accountRequest.User_id); err == model.IDnotFound || err != nil {
		return nil, err
	}

	if _, err := create.clientDatabase.Get(&accountRequest.Client_id); err == model.IDnotFound || err != nil {
		return nil, err
	}

	accountID, err := create.accountDatabase.Create(accountRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Account created")
	return create.getService.Account(*accountID)
}

func (create createImpl) CreateClient(clientRequest *client.ClientRequest) (*client.ClientResponse, *model.Erro) {
	if clientRequest.User_id == "" || clientRequest.Document == "" || clientRequest.Name == "" {
		log.Warn().Msg("Missing credentials on creating client")
		return nil, &model.Erro{Err: errors.New("Missing credentials for creating client"), HttpCode: http.StatusBadRequest}
	}

	if _, err := create.userDatabase.Get(&clientRequest.User_id); err == model.IDnotFound || err != nil {
		return nil, err
	}
	// verify user id exists, PERMISSION MUST BE of user to create
	log.Info().Msg("Client created: " + clientRequest.Client_id)
	clientID, err := create.clientDatabase.Create(clientRequest)
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
