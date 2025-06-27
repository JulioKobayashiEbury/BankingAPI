package service

import (
	"errors"
	"net/http"

	"BankingAPI/internal/gateway"
	"BankingAPI/internal/model"

	"BankingAPI/internal/model/transfer"

	"github.com/rs/zerolog/log"
)

type transferImpl struct {
	transferDatabase model.RepositoryInterface
	accountService   AccountService
	userService      UserService
	gateway          gateway.Gateway
}

func NewTransferService(transferDB model.RepositoryInterface, accountServe AccountService, userServe UserService, extGateway gateway.Gateway) TransferService {
	return transferImpl{
		transferDatabase: transferDB,
		userService:      userServe,
		accountService:   accountServe,
		gateway:          extGateway,
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
	if _, err := service.userService.Get(&transferRequest.User_to); err != nil {
		return nil, &model.Erro{Err: errors.New("User for account to was not found!"), HttpCode: http.StatusBadRequest}
	}

	accountTo, err := service.accountService.Get(&transferRequest.Account_to)
	if err == model.IDnotFound || err != nil {
		return nil, err
	}
	accountFrom, err := service.accountService.Get(&transferRequest.Account_id)
	if err == model.IDnotFound || err != nil {
		return nil, err
	}
	if accountFrom.Status != "active" || accountTo.Status != "active" {
		return nil, &model.Erro{Err: errors.New("one of the accounts is not active"), HttpCode: http.StatusBadRequest}
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

func (service transferImpl) ProcessExternalTransfer(transferRequest *transfer.Transfer) (*transfer.Transfer, *model.Erro) {
	_, err := service.accountService.Get(&transferRequest.Account_id)
	if err != nil {
		if err == model.IDnotFound {
			return service.OutsideToInside(transferRequest)
		}
		return nil, err
	}

	if _, err := service.accountService.Get(&transferRequest.Account_to); err != nil {
		if err == model.IDnotFound {
			return service.InsideToOutside(transferRequest)
		}
		return nil, err
	}

	return nil, &model.Erro{Err: errors.New("Neither account to or account from are from this system!"), HttpCode: http.StatusBadRequest}
}

func (service transferImpl) InsideToOutside(transferRequest *transfer.Transfer) (*transfer.Transfer, *model.Erro) {
	if err := service.gateway.Send(transferRequest); err != nil {
		return nil, err
	}

	accountFrom, err := service.accountService.Get(&transferRequest.Account_id)
	if err != nil {
		return nil, err
	}

	transferResponse, err := service.Create(transferRequest)
	if err != nil {
		return nil, err
	}

	accountFrom.Balance -= transferRequest.Value

	if _, err := service.accountService.Update(accountFrom); err != nil {
		if err := service.transferDatabase.Delete(&transferResponse.Transfer_id); err != nil {
			log.Error().Msg("could not delete trasnfer after account update failed " + transferResponse.Transfer_id)
			return nil, &model.Erro{Err: errors.New("could not delete trasnfer after account update failed"), HttpCode: http.StatusInternalServerError}
		}
		log.Error().Msg("Failed to complete transfer, update account failed!")
		return nil, &model.Erro{Err: errors.New("Failed to complete transfer, update account failed!"), HttpCode: http.StatusInternalServerError}
	}

	log.Info().Msg("External transfer is completed: " + transferResponse.Transfer_id)

	return transferResponse, nil
}

func (service transferImpl) OutsideToInside(transferRequest *transfer.Transfer) (*transfer.Transfer, *model.Erro) {
	accountTo, err := service.accountService.Get(&transferRequest.Account_to)
	if err != nil {
		return nil, err
	}

	transferResponse, err := service.Create(transferRequest)
	if err != nil {
		return nil, err
	}

	accountTo.Balance += transferRequest.Value

	if _, err := service.accountService.Update(accountTo); err != nil {
		if err := service.transferDatabase.Delete(&transferResponse.Transfer_id); err != nil {
			log.Error().Msg("could not delete trasnfer after account update failed " + transferResponse.Transfer_id)
			return nil, &model.Erro{Err: errors.New("could not delete trasnfer after account update failed"), HttpCode: http.StatusInternalServerError}
		}
		log.Error().Msg("Failed to complete transfer, update account failed!")
		return nil, &model.Erro{Err: errors.New("Failed to complete transfer, update account failed!"), HttpCode: http.StatusInternalServerError}
	}

	log.Info().Msg("External transfer is completed: " + transferResponse.Transfer_id)

	return transferResponse, nil
}
