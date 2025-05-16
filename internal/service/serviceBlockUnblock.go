package service

import (
	"fmt"

	"BankingAPI/internal/repository"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func AccountBlock(accountID uint32) error {
	if err := toggleStatus(false, &accountID, "accounts"); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func AccountUnBlock(accountID uint32) error {
	if err := toggleStatus(true, &accountID, "accounts"); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func ClientBlock(clientID uint32) error {
	if err := toggleStatus(false, &clientID, "clients"); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func ClientUnBlock(clientID uint32) error {
	if err := toggleStatus(true, &clientID, "clients"); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func UserBlock(userID uint32) error {
	if err := toggleStatus(false, &userID, "users"); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func UserUnBlock(userID uint32) error {
	if err := toggleStatus(true, &userID, "users"); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func toggleStatus(status bool, typesID *uint32, collection string) error {
	updates := []firestore.Update{
		{
			Path:  "status",
			Value: fmt.Sprintf("%v", status),
		},
	}
	// put account into db again
	if err := repository.UpdateTypesDB(&updates, typesID, &collection); err != nil {
		return err
	}
	return nil
}
