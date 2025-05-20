package service

import (
	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func AccountBlock(accountID string) *model.Erro {
	if err := toggleStatus(false, &accountID, repository.AccountsPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func AccountUnBlock(accountID string) *model.Erro {
	if err := toggleStatus(true, &accountID, repository.AccountsPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func ClientBlock(clientID string) *model.Erro {
	if err := toggleStatus(false, &clientID, repository.ClientPath); err != nil {
		return err
	}
	log.Info().Msg("Client Blocked")
	return nil
}

func ClientUnBlock(clientID string) *model.Erro {
	if err := toggleStatus(true, &clientID, repository.ClientPath); err != nil {
		return err
	}
	log.Info().Msg("Client UnBlocked")
	return nil
}

func UserBlock(userID string) *model.Erro {
	if err := toggleStatus(false, &userID, repository.UsersPath); err != nil {
		return err
	}
	log.Info().Msg("User Blocked")
	return nil
}

func UserUnBlock(userID string) *model.Erro {
	if err := toggleStatus(true, &userID, repository.UsersPath); err != nil {
		return err
	}
	log.Info().Msg("User UnBlocked")
	return nil
}

func toggleStatus(status bool, typesID *string, collection string) *model.Erro {
	updates := []firestore.Update{
		{
			Path:  "status",
			Value: status,
		},
	}
	// put account into db again
	if err := repository.UpdateTypesDB(&updates, typesID, collection); err != nil {
		return err
	}
	return nil
}
