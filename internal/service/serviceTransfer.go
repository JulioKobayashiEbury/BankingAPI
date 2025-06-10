package service

import (
	"BankingAPI/internal/model"

	"BankingAPI/internal/model/transfer"

	"github.com/rs/zerolog/log"
)

type transferImpl struct {
	transferDatabase model.RepositoryInterface
	accountService   AccountService
}

func NewTransferService(transferDB model.RepositoryInterface, accountServe AccountService) TransferService {
	return transferImpl{
		transferDatabase: transferDB,
		accountService:   accountServe,
	}
}

func (service transferImpl) Create(transferRequest *transfer.Transfer) (*transfer.Transfer, *model.Erro) {
	obj, err := service.transferDatabase.Create(transferRequest)
	if err != nil {
		return nil, err
	}
	transferResponse, ok := obj.(*transfer.Transfer)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return transferResponse, nil
}

func (service transferImpl) Delete(id *string) *model.Erro {
	if err := service.transferDatabase.Delete(id); err != nil {
		return err
	}
	return nil
}

func (service transferImpl) Get(id *string) (*transfer.Transfer, *model.Erro) {
	obj, err := service.transferDatabase.Get(id)
	if err != nil {
		return nil, err
	}

	transferResposne, ok := obj.(*transfer.Transfer)
	if !ok {
		return nil, model.DataTypeWrong
	}

	return transferResposne, nil
}

func (service transferImpl) GetAll() (*[]transfer.Transfer, *model.Erro) {
	obj, err := service.transferDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	transfers, ok := obj.(*[]transfer.Transfer)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return transfers, nil
}

func (service transferImpl) ProcessNewTransfer(transferRequest *transfer.Transfer) (*transfer.Transfer, *model.Erro) {
	accountTo, err := service.accountService.Get(&transferRequest.Account_to)
	if err == model.IDnotFound || err != nil {
		return nil, err
	}
	accountFrom, err := service.accountService.Get(&transferRequest.Account_id)
	if err == model.IDnotFound || err != nil {
		return nil, err
	}

	accountTo.Balance += transferRequest.Value
	accountFrom.Balance -= transferRequest.Value

	transferResponse, err := service.Create(transferRequest)
	if err != nil {
		return nil, err
	}
	if _, err := service.accountService.Update(accountTo); err != nil {
		if err := service.Delete(&transferResponse.Transfer_id); err != nil {
			return nil, err
		}
		log.Error().Msg("Update Account Receiving Transfer failed, transfer canceled")
		return nil, err
	}

	if _, err := service.accountService.Update(accountFrom); err != nil {
		accountTo.Balance -= transferRequest.Value
		if _, err := service.accountService.Update(accountTo); err != nil {
			return nil, err
		}
		if err := service.Delete(&transferResponse.Transfer_id); err != nil {
			return nil, err
		}
		log.Error().Msg("Update Account Sending Transfer failed, transfer canceled")
	}
	log.Info().Msg("Transfer was succesful: " + transferRequest.Transfer_id + " to " + transferRequest.Account_to)
	return transferResponse, nil
}
