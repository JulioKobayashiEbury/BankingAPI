package client

import (
	"context"

	"BankingAPI/internal/model"
)

var singleton *ClientRepository

type MockClientRepository struct {
	ClientMap *map[string]Client
}

func NewMockClientReposiory() ClientRepository {
	if singleton != nil {
		return *singleton
	}
	clientMap := make(map[string]Client)
	*singleton = MockClientRepository{
		ClientMap: &clientMap,
	}
	return *singleton
}

func (db MockClientRepository) Create(ctx context.Context, request *Client) (*Client, *model.Erro) {
	/*
		for {
			clientID := randomstring.String(10)
			if _, ok := (*db.ClientMap)[clientID]; !ok {
				clientRequest.Client_id = clientID
				(*db.ClientMap)[clientRequest.Client_id] = *clientRequest
				break
			}
		}
	*/
	(*db.ClientMap)[request.Client_id] = *request
	return request, nil
}

func (db MockClientRepository) Delete(ctx context.Context, id *string) *model.Erro {
	if _, ok := (*db.ClientMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.ClientMap, *id)
	return nil
}

func (db MockClientRepository) Get(ctx context.Context, id *string) (*Client, *model.Erro) {
	if clientResponse, ok := (*db.ClientMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &clientResponse, nil
	}
}

func (db MockClientRepository) Update(ctx context.Context, request *Client) *model.Erro {
	if _, ok := (*db.ClientMap)[request.Client_id]; !ok {
		return model.IDnotFound
	}
	return nil
}

func (db MockClientRepository) GetAll(ctx context.Context) (*[]Client, *model.Erro) {
	if len(*db.ClientMap) == 0 {
		return nil, model.IDnotFound
	}
	clients := make([]Client, 0, len(*db.ClientMap))
	for _, client := range *db.ClientMap {
		clients = append(clients, client)
	}
	return &clients, nil
}

func (db MockClientRepository) GetFilteredByID(ctx context.Context, filters *string) (*[]Client, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	clientSlice := make([]Client, 0, len(*db.ClientMap))
	/* for _, clientResponse := range *db.ClientMap {
		match := true
		for _, filter := range *filters {
			token := model.TokenizeFilters(&filter)
			if len(*token) != 3 {
				return nil, model.InvalidFilterFormat
			}
			field := (*token)[0]
			operator := (*token)[1]
			value := (*token)[2]

			switch field {
			case "user_id":
				if operator == "==" && clientResponse.User_id != value {
					match = false
				}
			// Add more fields as necessary
			default:
				match = false
			}
		}
		if match {
			clientSlice = append(clientSlice, clientResponse)
		}
	}
	*/
	return &clientSlice, nil
}
