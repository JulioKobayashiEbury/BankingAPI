package service

import (
	"context"
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type clientServiceImpl struct {
	clientDatabase  client.ClientRepository
	userService     UserService
	accountDatabase account.AccountRepository
}

func NewClientService(clientDB client.ClientRepository, userServe UserService, accountDB account.AccountRepository) ClientService {
	return clientServiceImpl{
		clientDatabase:  clientDB,
		userService:     userServe,
		accountDatabase: accountDB,
	}
}

func (service clientServiceImpl) Create(ctx context.Context, clientRequest *client.Client) (*client.Client, *echo.HTTPError) {
	if clientRequest.User_id == "" || clientRequest.Document == "" || clientRequest.Name == "" {
		log.Warn().Msg("Missing credentials on creating client")
		return nil, model.ErrMissingCredentials
	}

	if _, err := service.userService.Get(ctx, &clientRequest.User_id); err == model.ErrIDnotFound || err != nil {
		return nil, err
	}
	// verify user id exists, PERMISSION MUST BE of user to create
	clientResponse, err := service.clientDatabase.Create(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Client created: " + clientResponse.Client_id)
	return clientResponse, nil
}

func (service clientServiceImpl) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if err := service.clientDatabase.Delete(ctx, id); err != nil {
		return err
	}
	log.Info().Msg("Client deleted: " + *id)
	return nil
}

func (service clientServiceImpl) Get(ctx context.Context, id *string) (*client.Client, *echo.HTTPError) {
	client, err := service.clientDatabase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Client returned: " + *id)
	return client, nil
}

func (service clientServiceImpl) Update(ctx context.Context, clientRequest *client.Client) (*client.Client, *echo.HTTPError) {
	clientResponse, err := service.Get(ctx, &clientRequest.Client_id)
	if err != nil {
		return nil, err
	}

	if clientRequest.User_id != "" {
		clientResponse.User_id = clientRequest.User_id
	}
	if clientRequest.Name != "" {
		clientResponse.Name = clientRequest.Name
	}
	if clientRequest.Document != "" {
		clientResponse.Document = clientRequest.Document
	}
	// monta struct de update
	if err := service.clientDatabase.Update(ctx, clientResponse); err != nil {
		return nil, err
	}
	log.Info().Msg("Update was succesful (client): " + clientResponse.Client_id)
	return service.Get(ctx, &clientResponse.Client_id)
}

func (service clientServiceImpl) GetAll(ctx context.Context) (*[]client.Client, *echo.HTTPError) {
	clients, err := service.clientDatabase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (service clientServiceImpl) Report(ctx context.Context, clientId *string) (*client.ClientReport, *echo.HTTPError) {
	clientInfo, err := service.Get(ctx, clientId)
	if err != nil {
		return nil, err
	}
	accounts, err := service.accountDatabase.GetFilteredByClientID(ctx, clientId)
	if err != nil {
		return nil, err
	}
	return &client.ClientReport{
		Client_id:     clientInfo.Client_id,
		User_id:       clientInfo.User_id,
		Name:          clientInfo.Name,
		Document:      clientInfo.Document,
		Register_date: clientInfo.Register_date,
		Accounts:      accounts,
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}
