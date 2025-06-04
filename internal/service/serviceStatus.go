package service

import (
	"BankingAPI/internal/model"

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
	getService     ServiceGet
}

func NewStatusService(userDB model.RepositoryInterface, clientDB model.RepositoryInterface, accountDB model.RepositoryInterface, get ServiceGet) StatusService {
	return statusImpl{
		userDatabase:   userDB,
		clientDatabase: clientDB,
		accountDatabse: accountDB,
		getService:     get,
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
	account, err := si.getService.Account(*accountID)
	if err != nil {
		return err
	}
	account.Status = status
	if err := si.accountDatabse.Update(account); err != nil {
		return err
	}
	return nil
}

func (si statusImpl) toggleClientStatus(status bool, clientID *string) *model.Erro {
	client, err := si.getService.Client(*clientID)
	if err != nil {
		return err
	}
	client.Status = status
	if err := si.clientDatabase.Update(client); err != nil {
		return err
	}
	return nil
}

func (si statusImpl) toggleUserStatus(status bool, userID *string) *model.Erro {
	user, err := si.getService.User(*userID)
	if err != nil {
		return err
	}
	user.Status = status
	if err := si.userDatabase.Update(userID); err != nil {
		return err
	}
	return nil
}
