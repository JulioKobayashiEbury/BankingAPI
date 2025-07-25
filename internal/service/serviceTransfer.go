package service

import (
	"context"
	"errors"
	"net/http"

	"BankingAPI/internal/gateway"
	"BankingAPI/internal/model"

	"BankingAPI/internal/model/transfer"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type transferImpl struct {
	transferDatabase transfer.TransferRepository
	accountService   AccountService
	userService      UserService
	gateway          gateway.Gateway
}

func NewTransferService(transferDB transfer.TransferRepository, accountServe AccountService, userServe UserService, extGateway gateway.Gateway) TransferService {
	return transferImpl{
		transferDatabase: transferDB,
		userService:      userServe,
		accountService:   accountServe,
		gateway:          extGateway,
	}
}

func (service transferImpl) Create(ctx context.Context, transferRequest *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError) {
	transferResponse, err := service.transferDatabase.Create(ctx, transferRequest)
	if err != nil {
		return nil, err
	}
	return transferResponse, nil
}

func (service transferImpl) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if err := service.transferDatabase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (service transferImpl) Get(ctx context.Context, id *string) (*transfer.Transfer, *echo.HTTPError) {
	transferResponse, err := service.transferDatabase.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return transferResponse, nil
}

func (service transferImpl) GetAll(ctx context.Context) (*[]transfer.Transfer, *echo.HTTPError) {
	transfers, err := service.transferDatabase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (service transferImpl) ProcessNewTransfer(ctx context.Context, transferRequest *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError) {
	if _, err := service.userService.Get(ctx, &transferRequest.User_to); err != nil {
		return nil, &echo.HTTPError{Internal: errors.New("user for account to was not found!"), Code: http.StatusBadRequest, Message: "user for account to was not found!"}
	}

	accountTo, err := service.accountService.Get(ctx, &transferRequest.Account_to)
	if err == model.ErrIDnotFound || err != nil {
		return nil, err
	}
	accountFrom, err := service.accountService.Get(ctx, &transferRequest.Account_id)
	if err == model.ErrIDnotFound || err != nil {
		return nil, err
	}
	if accountFrom.Status != "active" || accountTo.Status != "active" {
		return nil, &echo.HTTPError{Internal: errors.New("one of the accounts is not active"), Code: http.StatusBadRequest, Message: "one of the accounts is not active"}
	}

	accountTo.Balance += transferRequest.Value
	accountFrom.Balance -= transferRequest.Value

	transferResponse, err := service.Create(ctx, transferRequest)
	if err != nil {
		return nil, err
	}
	if _, err := service.accountService.Update(ctx, accountTo); err != nil {
		if err := service.Delete(ctx, &transferResponse.Transfer_id); err != nil {
			return nil, err
		}
		log.Error().Msg("Update Account Receiving Transfer failed, transfer canceled")
		return nil, err
	}

	if _, err := service.accountService.Update(ctx, accountFrom); err != nil {
		accountTo.Balance -= transferRequest.Value
		if _, err := service.accountService.Update(ctx, accountTo); err != nil {
			return nil, err
		}
		if err := service.Delete(ctx, &transferResponse.Transfer_id); err != nil {
			return nil, err
		}
		log.Error().Msg("Update Account Sending Transfer failed, transfer canceled")
	}
	log.Info().Msg("Transfer was succesful: " + transferRequest.Transfer_id + " to " + transferRequest.Account_to)
	return transferResponse, nil
}

func (service transferImpl) ProcessExternalTransfer(ctx context.Context, transferRequest *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError) {
	_, err := service.accountService.Get(ctx, &transferRequest.Account_id)
	if err != nil {
		if err == model.ErrIDnotFound {
			return service.OutsideToInside(ctx, transferRequest)
		}
		return nil, err
	}

	if _, err := service.accountService.Get(ctx, &transferRequest.Account_to); err != nil {
		if err == model.ErrIDnotFound {
			return service.InsideToOutside(ctx, transferRequest)
		}
		return nil, err
	}

	return nil, &echo.HTTPError{Internal: errors.New("Neither account to or account from are from this system!"), Code: http.StatusBadRequest, Message: "Neither account to or account from are from this system!"}
}

func (service transferImpl) InsideToOutside(ctx context.Context, transferRequest *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError) {
	if err := service.gateway.Send(transferRequest); err != nil {
		return nil, err
	}

	accountFrom, err := service.accountService.Get(ctx, &transferRequest.Account_id)
	if err != nil {
		return nil, err
	}

	transferResponse, err := service.Create(ctx, transferRequest)
	if err != nil {
		return nil, err
	}

	accountFrom.Balance -= transferRequest.Value

	if _, err := service.accountService.Update(ctx, accountFrom); err != nil {
		if err := service.transferDatabase.Delete(ctx, &transferResponse.Transfer_id); err != nil {
			log.Error().Msg("could not delete trasnfer after account update failed " + transferResponse.Transfer_id)
			return nil, &echo.HTTPError{Internal: errors.New("could not delete trasnfer after account update failed"), Code: http.StatusInternalServerError, Message: "could not delete trasnfer after account update failed"}
		}
		log.Error().Msg("Failed to complete transfer, update account failed!")
		return nil, &echo.HTTPError{Internal: errors.New("Failed to complete transfer, update account failed!"), Code: http.StatusInternalServerError, Message: "Failed to complete transfer, update account failed!"}
	}

	log.Info().Msg("External transfer is completed: " + transferResponse.Transfer_id)

	return transferResponse, nil
}

func (service transferImpl) OutsideToInside(ctx context.Context, transferRequest *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError) {
	accountTo, err := service.accountService.Get(ctx, &transferRequest.Account_to)
	if err != nil {
		return nil, err
	}

	transferResponse, err := service.Create(ctx, transferRequest)
	if err != nil {
		return nil, err
	}

	accountTo.Balance += transferRequest.Value

	if _, err := service.accountService.Update(ctx, accountTo); err != nil {
		if err := service.transferDatabase.Delete(ctx, &transferResponse.Transfer_id); err != nil {
			log.Error().Msg("could not delete trasnfer after account update failed " + transferResponse.Transfer_id)
			return nil, &echo.HTTPError{Internal: errors.New("could not delete trasnfer after account update failed"), Code: http.StatusInternalServerError, Message: "could not delete trasnfer after account update failed"}
		}
		log.Error().Msg("Failed to complete transfer, update account failed!")
		return nil, &echo.HTTPError{Internal: errors.New("Failed to complete transfer, update account failed!"), Code: http.StatusInternalServerError, Message: "Failed to complete transfer, update account failed!"}
	}

	log.Info().Msg("External transfer is completed: " + transferResponse.Transfer_id)

	return transferResponse, nil
}
