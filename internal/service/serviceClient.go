package service

import (
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/client"

	"github.com/rs/zerolog/log"
)

type clientServiceImpl struct {
	clientDatabase     model.RepositoryInterface
	userService        UserService
	getFilteredService GetFilteredService
}

func NewClientService(clientDB model.RepositoryInterface, userServe UserService, getFilteredServe GetFilteredService) ClientService {
	return clientServiceImpl{
		clientDatabase:     clientDB,
		userService:        userServe,
		getFilteredService: getFilteredServe,
	}
}

func (service clientServiceImpl) Create(clientRequest *client.Client) (*client.Client, *model.Erro) {
	if clientRequest.User_id == "" || clientRequest.Document == "" || clientRequest.Name == "" {
		log.Warn().Msg("Missing credentials on creating client")
		return nil, ErrorMissingCredentials
	}

	if _, err := service.userService.Get(&clientRequest.User_id); err == model.IDnotFound || err != nil {
		return nil, err
	}
	// verify user id exists, PERMISSION MUST BE of user to create
	obj, err := service.clientDatabase.Create(clientRequest)
	if err != nil {
		return nil, err
	}
	clientResponse, ok := obj.(*client.Client)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Client created: " + clientResponse.Client_id)
	return clientResponse, nil
}

func (service clientServiceImpl) Delete(id *string) *model.Erro {
	if err := service.clientDatabase.Delete(id); err != nil {
		return err
	}
	log.Info().Msg("Client deleted: " + *id)
	return nil
}

func (service clientServiceImpl) Get(id *string) (*client.Client, *model.Erro) {
	obj, err := service.clientDatabase.Get(id)
	if err != nil {
		return nil, err
	}
	client, ok := obj.(*client.Client)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Client returned: " + *id)
	return client, nil
}

func (service clientServiceImpl) Update(clientRequest *client.Client) (*client.Client, *model.Erro) {
	client, err := service.Get(&clientRequest.Client_id)
	if err != nil {
		return nil, err
	}

	if clientRequest.User_id != "" {
		client.User_id = clientRequest.User_id
	}
	if clientRequest.Name != "" {
		client.Name = clientRequest.Name
	}
	if clientRequest.Document != "" {
		client.Document = clientRequest.Document
	}
	// monta struct de update
	if err := service.clientDatabase.Update(client); err != nil {
		return nil, err
	}
	log.Info().Msg("Update was succesful (client): " + client.Client_id)
	return service.Get(&client.Client_id)
}

func (service clientServiceImpl) GetAll() (*[]client.Client, *model.Erro) {
	obj, err := service.clientDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	clients, ok := obj.(*[]client.Client)
	if !ok {
		return nil, model.DataTypeWrong
	}

	return clients, nil
}

func (service clientServiceImpl) Status(id *string, status bool) *model.Erro {
	client, err := service.Get(id)
	if err != nil {
		return err
	}
	client.Status = status
	if err := service.clientDatabase.Update(client); err != nil {
		return err
	}
	return nil
}

func (service clientServiceImpl) Report(id *string) (*client.ClientReport, *model.Erro) {
	clientInfo, err := service.Get(id)
	if err != nil {
		return nil, err
	}
	accounts, err := service.getFilteredService.GetAccountsByClientID(id)
	if err != nil {
		return nil, err
	}
	return &client.ClientReport{
		Client_id:     clientInfo.Client_id,
		User_id:       clientInfo.User_id,
		Name:          clientInfo.Name,
		Document:      clientInfo.Document,
		Register_date: clientInfo.Register_date,
		Status:        clientInfo.Status,
		Accounts:      (*accounts),
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}
