package service

import (
	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

func AccountDelete(accountID string) *model.Erro {
	database := &account.AccountFirestore{}
	database.Request = &account.AccountRequest{
		Account_id: accountID,
	}
	if err := database.Delete(); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + accountID)
	return nil
}

func ClientDelete(clientID string) *model.Erro {
	database := &client.ClientFirestore{}
	database.Request = &client.ClientRequest{
		Client_id: clientID,
	}
	if err := database.Delete(); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + clientID)
	return nil
}

func UserDelete(userID string) *model.Erro {
	database := &user.UserFirestore{}
	database.Request = &user.UserRequest{
		User_id: userID,
	}
	if err := database.Delete(); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + userID)
	return nil
}
