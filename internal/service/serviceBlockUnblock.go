package service

import (
	"fmt"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func AccountBlock(accountID string) error {
	if err := toggleStatus(false, &accountID, model.AccountsPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func AccountUnBlock(accountID string) error {
	if err := toggleStatus(true, &accountID, model.AccountsPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func ClientBlock(clientID string) error {
	if err := toggleStatus(false, &clientID, model.ClientPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func ClientUnBlock(clientID string) error {
	if err := toggleStatus(true, &clientID, model.ClientPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func UserBlock(userID string) error {
	if err := toggleStatus(false, &userID, model.UsersPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func UserUnBlock(userID string) error {
	if err := toggleStatus(true, &userID, model.UsersPath); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func toggleStatus(status bool, typesID *string, collection string) error {
	updates := []firestore.Update{
		{
			Path:  "status",
			Value: fmt.Sprintf("%v", status),
		},
	}
	// put account into db again
	if err := model.UpdateTypesDB(&updates, typesID, collection); err != nil {
		return err
	}
	return nil
}
