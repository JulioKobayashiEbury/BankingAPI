package service

import (
	"BankingAPI/internal/model"

	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/transfer"

	"github.com/rs/zerolog/log"
)

type transferImpl struct {
	accountDatabase  model.RepositoryInterface
	transferDatabase model.RepositoryInterface
}

func NewTransferService(accountDB model.RepositoryInterface, transferDB model.RepositoryInterface) TransferService {
	return transferImpl{
		accountDatabase:  accountDB,
		transferDatabase: transferDB,
	}
}

func (transfer transferImpl) Create(*transfer.Transfer) (*string, *model.Erro)
func (transfer transferImpl) Delete(*string) *model.Erro
func (transfer transferImpl) GetAll(*string) ([]*transfer.Transfer, *model.Erro)

func (transfer transferImpl) ProcessNewTransfer(transferRequest *transfer.Transfer) (*string, *model.Erro) {
	obj, err := transfer.accountDatabase.Get(&transferRequest.Account_to)
	if err == model.IDnotFound || err != nil {
		return nil, err
	}
	accountTo, ok := obj.(*account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	obj, err = transfer.accountDatabase.Get(&transferRequest.Account_id)
	if err == model.IDnotFound || err != nil {
		return nil, err
	}
	accountFrom, ok := obj.(*account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}

	accountTo.Balance += transferRequest.Value
	accountFrom.Balance -= transferRequest.Value

	transferID, err := transfer.transferDatabase.Create(transferRequest)
	if err != nil {
		return nil, err
	}
	if err := transfer.accountDatabase.Update(accountTo); err != nil {
		if err := transfer.transferDatabase.Delete(transferID); err != nil {
			return nil, err
		}
		log.Error().Msg("Update Account Receiving Transfer failed, transfer canceled")
		return nil, err
	}

	if err := transfer.accountDatabase.Update(accountFrom); err != nil {
		accountTo.Balance -= transferRequest.Value
		if err := transfer.accountDatabase.Update(accountTo); err != nil {
			return nil, err
		}
		if err := transfer.transferDatabase.Delete(transferID); err != nil {
			return nil, err
		}
		log.Error().Msg("Update Account Sending Transfer failed, transfer canceled")
	}
	log.Info().Msg("Transfer was succesful: " + *transferID + " to " + transferRequest.Account_to)
	return transferID, nil
}
