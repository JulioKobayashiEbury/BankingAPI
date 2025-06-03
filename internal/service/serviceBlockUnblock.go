package service

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

type StatusService interface {
	AccountBlock(accountID string) *model.Erro
	AccountUnBlock(accountID string) *model.Erro
	ClientBlock(clientID string) *model.Erro
	ClientUnBlock(clientID string) *model.Erro
	UserBlock(userID string) *model.Erro
	UserUnBlock(userID string) *model.Erro
}

type statusImpl struct {
	userDatabase   model.RepositoryInterface
	clientDatabase model.RepositoryInterface
	accountDatabse model.RepositoryInterface
}

func NewStatusService(dbClient *firestore.Client) StatusService {
	return statusImpl{
		userDatabase:   user.NewUserFireStore(dbClient),
		clientDatabase: client.NewClientFirestore(dbClient),
		accountDatabse: client.NewClientFirestore(dbClient),
	}
}

func (si statusImpl) AccountBlock(accountID string) *model.Erro {
	if err := si.toggleAccountStatus(false, &accountID); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func (si statusImpl) AccountUnBlock(accountID string) *model.Erro {
	if err := si.toggleAccountStatus(true, &accountID); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func (si statusImpl) ClientBlock(clientID string) *model.Erro {
	if err := si.toggleClientStatus(false, &clientID); err != nil {
		return err
	}
	log.Info().Msg("Client Blocked")
	return nil
}

func (si statusImpl) ClientUnBlock(clientID string) *model.Erro {
	if err := si.toggleClientStatus(true, &clientID); err != nil {
		return err
	}
	log.Info().Msg("Client UnBlocked")
	return nil
}

func (si statusImpl) UserBlock(userID string) *model.Erro {
	if err := si.toggleUserStatus(false, &userID); err != nil {
		return err
	}
	log.Info().Msg("User Blocked")
	return nil
}

func (si statusImpl) UserUnBlock(userID string) *model.Erro {
	if err := si.toggleUserStatus(true, &userID); err != nil {
		return err
	}
	log.Info().Msg("User UnBlocked")
	return nil
}

func (si statusImpl) toggleAccountStatus(status bool, accountID *string) *model.Erro {
	si.accountDatabse.AddUpdate("status", status)
	if err := si.accountDatabse.Update(accountID); err != nil {
		return err
	}
	return nil
}

func (si statusImpl) toggleClientStatus(status bool, clientID *string) *model.Erro {
	si.clientDatabase.AddUpdate("status", status)
	if err := si.clientDatabase.Update(clientID); err != nil {
		return err
	}
	return nil
}

func (si statusImpl) toggleUserStatus(status bool, userID *string) *model.Erro {
	si.userDatabase.AddUpdate("status", status)
	if err := si.userDatabase.Update(userID); err != nil {
		return err
	}
	return nil
}
