package service

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

func AccountBlock(accountID string) *model.Erro {
	if err := toggleAccountStatus(false, &accountID); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func AccountUnBlock(accountID string) *model.Erro {
	if err := toggleAccountStatus(true, &accountID); err != nil {
		return err
	}
	log.Info().Msg("Account Blocked")
	return nil
}

func ClientBlock(clientID string) *model.Erro {
	if err := toggleClientStatus(false, &clientID); err != nil {
		return err
	}
	log.Info().Msg("Client Blocked")
	return nil
}

func ClientUnBlock(clientID string) *model.Erro {
	if err := toggleClientStatus(true, &clientID); err != nil {
		return err
	}
	log.Info().Msg("Client UnBlocked")
	return nil
}

func UserBlock(userID string) *model.Erro {
	if err := toggleUserStatus(false, &userID); err != nil {
		return err
	}
	log.Info().Msg("User Blocked")
	return nil
}

func UserUnBlock(userID string) *model.Erro {
	if err := toggleUserStatus(true, &userID); err != nil {
		return err
	}
	log.Info().Msg("User UnBlocked")
	return nil
}

func toggleAccountStatus(status bool, accountID *string) *model.Erro {
	database := &account.AccountFirestore{}
	database.Request = &account.AccountRequest{
		Account_id: *accountID,
	}
	database.AddUpdate("status", status)
	if err := database.Update(); err != nil {
		return err
	}
	return nil
}

func toggleClientStatus(status bool, clientID *string) *model.Erro {
	database := &client.ClientFirestore{}
	database.Request = &client.ClientRequest{
		Client_id: *clientID,
	}
	database.AddUpdate("status", status)
	if err := database.Update(); err != nil {
		return err
	}
	return nil
}

func toggleUserStatus(status bool, userID *string) *model.Erro {
	database := &user.UserFirestore{}
	database.Request = &user.UserRequest{
		User_id: *userID,
	}
	database.AddUpdate("status", status)
	if err := database.Update(); err != nil {
		return err
	}
	return nil
}
