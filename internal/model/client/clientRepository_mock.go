package client

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo/v4"
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

func (db MockClientRepository) Create(ctx context.Context, request *Client) (*Client, *echo.HTTPError) {
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

func (db MockClientRepository) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if _, ok := (*db.ClientMap)[*id]; !ok {
		return model.ErrIDnotFound
	}
	delete(*db.ClientMap, *id)
	return nil
}

func (db MockClientRepository) Get(ctx context.Context, id *string) (*Client, *echo.HTTPError) {
	if clientResponse, ok := (*db.ClientMap)[*id]; !ok {
		return nil, model.ErrIDnotFound
	} else {
		return &clientResponse, nil
	}
}

func (db MockClientRepository) Update(ctx context.Context, request *Client) *echo.HTTPError {
	if _, ok := (*db.ClientMap)[request.Client_id]; !ok {
		return model.ErrIDnotFound
	}
	return nil
}

func (db MockClientRepository) GetAll(ctx context.Context) (*[]Client, *echo.HTTPError) {
	if len(*db.ClientMap) == 0 {
		return nil, model.ErrIDnotFound
	}
	clients := make([]Client, 0, len(*db.ClientMap))
	for _, client := range *db.ClientMap {
		clients = append(clients, client)
	}
	return &clients, nil
}

func (db MockClientRepository) GetFilteredByUserID(ctx context.Context, userID *string) (*[]Client, *echo.HTTPError) {
	if userID == nil || len(*userID) == 0 {
		return nil, model.ErrFilterNotSet
	}
	clientSlice := make([]Client, 0, len(*db.ClientMap))

	for _, client := range *db.ClientMap {
		if client.User_id == *userID {
			clientSlice = append(clientSlice, client)
		}
	}
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
