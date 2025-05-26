package service

import (
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

func Account(accountID string) (*model.AccountResponse, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&accountID, repository.AccountsPath)
	if err != nil {
		return nil, err
	}
	var accountResponse model.AccountResponse
	if err := docSnapshot.DataTo(&accountResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponse.Account_id = accountID

	log.Info().Msg("Account returned: " + accountID)
	return &accountResponse, nil
}

func GetAccountByFilterAndOrder(listRequest *model.ListRequest) (*[]model.AccountResponse, *model.Erro) {
	return nil, nil
}

func Client(clientID string) (*model.ClientResponse, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&clientID, repository.ClientPath)
	if err != nil {
		return nil, err
	}
	var clientResponse model.ClientResponse
	if err := docSnapshot.DataTo(&clientResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	clientResponse.Client_id = clientID

	log.Info().Msg("Client returned: " + clientID)
	return &clientResponse, nil
}

func User(userID string) (*model.UserResponse, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&userID, repository.UsersPath)
	if err != nil {
		return nil, err
	}

	var userResponse model.UserResponse
	if err := docSnapshot.DataTo(&userResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	userResponse.User_id = userID

	log.Info().Msg("User returned: " + userID)
	return &userResponse, nil
}

// implement getallby(each type)
func GetAllTransfers(accountID *string) (*[]model.TransferResponse, *model.Erro) {
	docSnapshots, err := repository.GetAllByTypeDB(repository.TransfersPath)
	if err != nil {
		return nil, err
	}
	tranfersReponseSlice := make([]model.TransferResponse, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		transferResponse := model.TransferResponse{}
		if err := docSnap.DataTo(&transferResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		transferResponse.Transfer_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		if *accountID == transferResponse.Account_id {
			tranfersReponseSlice = append(tranfersReponseSlice, transferResponse)
		}
	}

	return &tranfersReponseSlice, nil
}

func GetAllAutoDebits(accountID *string) (*[]model.AutomaticDebitResponse, *model.Erro) {
	docSnapshots, err := repository.GetAllByTypeDB(repository.AutoDebit)
	if err != nil {
		return nil, err
	}
	autoDebitsResponseSlice := make([]model.AutomaticDebitResponse, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		autoDebitResponse := model.AutomaticDebitResponse{}
		if err := docSnap.DataTo(&autoDebitResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		autoDebitResponse.Debit_id = docSnap.Ref.ID
		if *accountID == autoDebitResponse.Account_id {
			autoDebitsResponseSlice = append(autoDebitsResponseSlice, autoDebitResponse)
		}
	}
	return &autoDebitsResponseSlice, nil
}

func GetAllWithdrawals(accountID *string) (*[]model.WithdrawalResponse, *model.Erro) {
	docSnapshots, err := repository.GetAllByTypeDB(repository.WithdrawalsPath)
	if err != nil {
		return nil, err
	}
	withdrawalsResponseSlice := make([]model.WithdrawalResponse, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		withdrawalsResponse := model.WithdrawalResponse{}
		if err := docSnap.DataTo(&withdrawalsResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		withdrawalsResponse.Withdrawal_id = docSnap.Ref.ID
		if *accountID == withdrawalsResponse.Withdrawal_id {
			withdrawalsResponseSlice = append(withdrawalsResponseSlice, withdrawalsResponse)
		}
	}
	return &withdrawalsResponseSlice, nil
}

func GetAllDeposits(accountID *string) (*[]model.DepositResponse, *model.Erro) {
	docSnapshots, err := repository.GetAllByTypeDB(repository.DepositPath)
	if err != nil {
		return nil, err
	}
	depositsResponseSlice := make([]model.DepositResponse, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		depositResponse := model.DepositResponse{}
		if err := docSnap.DataTo(&depositResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		depositResponse.Deposit_id = docSnap.Ref.ID
		if *accountID == depositResponse.Account_id {
			depositsResponseSlice = append(depositsResponseSlice, depositResponse)
		}
	}
	return &depositsResponseSlice, nil
}
