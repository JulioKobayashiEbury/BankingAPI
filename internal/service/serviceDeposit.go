package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/deposit"

	"github.com/rs/zerolog/log"
)

type depositImpl struct {
	depositDatabase model.RepositoryInterface
	accountDatabase model.RepositoryInterface
	getService      ServiceGet
}

func NewDepositService(depositDB model.RepositoryInterface, accountDB model.RepositoryInterface, get ServiceGet) DepositService {
	return depositImpl{
		depositDatabase: depositDB,
		accountDatabase: accountDB,
		getService:      get,
	}
}

func (deposit depositImpl) Create(*deposit.Deposit) (*string, *model.Erro)
func (deposit depositImpl) Delete(*string) *model.Erro
func (deposit depositImpl) GetAll(*string) ([]*deposit.Deposit, *model.Erro)

func (deposit depositImpl) ProcessDeposit(depositRequest *deposit.Deposit) (*string, *model.Erro) {
	accountRequest, err := deposit.getService.Account(depositRequest.Account_id)
	if err != nil {
		return nil, err
	}
	if ok, err := verifyDeposit(depositRequest, accountRequest); !ok {
		return nil, err
	}
	accountRequest.Balance = accountRequest.Balance + depositRequest.Deposit

	depositID, err := deposit.depositDatabase.Create(depositRequest)
	if err != nil {
		return nil, err
	}

	if err := deposit.accountDatabase.Update(accountRequest); err != nil {
		if err := deposit.depositDatabase.Delete(depositID); err != nil {
			log.Panic().Msg("Error deleting deposit after update failure: " + err.Err.Error())
		}
		return nil, err
	}
	log.Info().Msg("Deposit created: " + depositRequest.Account_id)
	return depositID, nil
}

func verifyDeposit(depositRequest *deposit.Deposit, accountResponse *account.Account) (bool, *model.Erro) {
	if accountResponse.Client_id != depositRequest.Client_id {
		return false, &model.Erro{Err: errors.New("Client ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountResponse.User_id != depositRequest.User_id {
		return false, &model.Erro{Err: errors.New("User ID not valid"), HttpCode: http.StatusBadRequest}
	}
	if accountResponse.Agency_id != depositRequest.Agency_id {
		return false, &model.Erro{Err: errors.New("Agency ID not valid"), HttpCode: http.StatusBadRequest}
	}
	return true, nil
}
