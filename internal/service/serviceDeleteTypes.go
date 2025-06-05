package service

import (
	model "BankingAPI/internal/model"

	"github.com/rs/zerolog/log"
)

type DeleteService interface {
	AccountDelete(accountID string) *model.Erro
	ClientDelete(clientID string) *model.Erro
	UserDelete(userID string) *model.Erro
}

type deleteImpl struct {
	userDatabase    model.RepositoryInterface
	clientDatabase  model.RepositoryInterface
	accountDatabase model.RepositoryInterface
}

func NewDeleteService(userDB model.RepositoryInterface, clientDB model.RepositoryInterface, accounDB model.RepositoryInterface) DeleteService {
	return deleteImpl{
		userDatabase:    userDB,
		clientDatabase:  clientDB,
		accountDatabase: accounDB,
	}
}

func (delete deleteImpl) ClientDelete(clientID string) *model.Erro {
	if err := delete.clientDatabase.Delete(&clientID); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + clientID)
	return nil
}

func (delete deleteImpl) UserDelete(userID string) *model.Erro {
	if err := delete.userDatabase.Delete(&userID); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + userID)
	return nil
}
