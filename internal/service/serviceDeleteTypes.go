package service

import (
	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

func AccountDelete(accountID string) *model.Erro {
	if err := repository.DeleteObject(&accountID, repository.AccountsPath); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + accountID)
	return nil
}

func ClientDelete(clientID string) *model.Erro {
	if err := repository.DeleteObject(&clientID, repository.ClientPath); err != nil {
		return err
	}
	log.Info().Msg("Client deleted: " + clientID)
	return nil
}

func UserDelete(userID string) *model.Erro {
	if err := repository.DeleteObject(&userID, repository.UsersPath); err != nil {
		return err
	}
	log.Info().Msg("User deleted: " + userID)
	return nil
}
