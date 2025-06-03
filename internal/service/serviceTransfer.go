package service

import (
	"BankingAPI/internal/model"

	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/transfer"

	"github.com/rs/zerolog/log"
)

type ServiceTransfer interface {
	ProcessNewTransfer(*transfer.TransferRequest) *model.Erro
}

type transferImpl struct {
	accountDatabase  account.AccountFirestore
	transferDatabase transfer.TransferFirestore
}

func ProcessNewTransfer(transferRequest *transfer.TransferRequest) *model.Erro {
	accountToDatabase := &account.AccountFirestore{}
	accountToDatabase.Request = &account.AccountRequest{
		Account_id: transferRequest.Account_to,
	}

	accountFromDatabase := &account.AccountFirestore{}
	accountFromDatabase.Request = &account.AccountRequest{
		Account_id: transferRequest.Account_id,
	}

	if err := accountToDatabase.Get(); err == model.IDnotFound || err != nil {
		return err
	}
	if err := accountFromDatabase.Get(); err == model.IDnotFound || err != nil {
		return err
	}

	accountToDatabase.Response.Balance += transferRequest.Value
	accountFromDatabase.Response.Balance -= transferRequest.Value

	transferDatabase := &transfer.TransferFirestore{
		Request: transferRequest,
	}
	if err := transferDatabase.Create(); err != nil {
		return err
	}
	transferDatabase.Request.Transfer_id = transferDatabase.Response.Transfer_id
	accountToDatabase.AddUpdate("balance", accountToDatabase.Response.Balance)
	if err := accountToDatabase.Update(); err != nil {
		if err := transferDatabase.Delete(); err != nil {
			return err
		}
		log.Error().Msg("Update Account Receiving Transfer failed, transfer canceled")
		return err
	}
	accountFromDatabase.AddUpdate("balance", accountFromDatabase.Response.Balance)
	if err := accountFromDatabase.Update(); err != nil {
		accountToDatabase.AddUpdate("balance", accountToDatabase.Response.Balance-transferRequest.Value)
		if err := accountToDatabase.Update(); err != nil {
			return err
		}
		if err := transferDatabase.Delete(); err != nil {
			return err
		}
		log.Error().Msg("Update Account Sending Transfer failed, transfer canceled")
	}
	log.Info().Msg("Transfer was succesful: " + transferRequest.Transfer_id + " to " + transferRequest.Account_to)
	return nil
}
