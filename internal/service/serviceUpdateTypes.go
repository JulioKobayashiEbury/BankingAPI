package service

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

type ServiceUpdate interface {
	UpdateAccount(accountRequest *account.Account) (*account.Account, *model.Erro)
	UpdateClient(Client *client.Client) (*client.Client, *model.Erro)
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

func (update updateImpl) UpdateClient(Client *client.Client) (*client.Client, *model.Erro) {
	Client, err := update.get.Client(Client.Client_id)
	if err != nil {
		return nil, err
	}

	if Client.User_id != "" {
		Client.User_id = Client.User_id
	}
	if Client.Name != "" {
		Client.Name = Client.Name
	}
	if Client.Document != "" {
		Client.Document = Client.Document
	}
	// monta struct de update
	if err := update.clientDatabase.Update(Client); err != nil {
		return nil, err
	}
	log.Info().Msg("Update was succesful (client): " + Client.Client_id)
	return update.get.Client(Client.Client_id)
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
