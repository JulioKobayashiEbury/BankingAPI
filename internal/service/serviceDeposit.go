package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/deposit"

	"github.com/rs/zerolog/log"
)

func ProcessDeposit(depositRequest *deposit.DepositRequest) *model.Erro {
	accountRequest, err := Account(depositRequest.Account_id)
	if err != nil {
		return err
	}
	if accountRequest.Client_id != depositRequest.Client_id {
		return &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountRequest.User_id != depositRequest.User_id {
		return &model.Erro{Err: errors.New("User ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountRequest.Agency_id != depositRequest.Agency_id {
		return &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}

	accountRequest.Balance = accountRequest.Balance + depositRequest.Deposit

	depositDatabase := &deposit.DepositFirestore{
		Request: depositRequest,
	}
	if err := depositDatabase.Create(); err != nil {
		return err
	}
	depositDatabase.Request.Deposit_id = depositDatabase.Response.Deposit_id

	accountDatabase := &account.AccountFirestore{}
	accountDatabase.Request.Account_id = depositRequest.Account_id
	accountDatabase.AddUpdate("balance", accountRequest.Balance)
	if err := accountDatabase.Update(); err != nil {
		if err := depositDatabase.Delete(); err != nil {
			log.Panic().Msg("Error deleting deposit after update failure: " + err.Err.Error())
		}
		return err
	}
	log.Info().Msg("Deposit created: " + depositRequest.Account_id)
	return nil
}
