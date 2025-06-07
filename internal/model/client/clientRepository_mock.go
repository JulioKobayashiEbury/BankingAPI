package client

import "BankingAPI/internal/model"

var singleton *model.RepositoryInterface

type MockClientRepository struct {
	ClientMap *map[string]Client
}

func NewMockClientReposiory() model.RepositoryInterface {
	if singleton != nil {
		return *singleton
	}
	clientMap := make(map[string]Client)
	*singleton = MockClientRepository{
		ClientMap: &clientMap,
	}
	return *singleton
}

func (db MockClientRepository) Create(request interface{}) (interface{}, *model.Erro) {
	clientRequest, ok := request.(*Client)
	if !ok {
		return nil, model.DataTypeWrong
	}
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
	(*db.ClientMap)[clientRequest.Client_id] = *clientRequest
	return &clientRequest, nil
}

func (db MockClientRepository) Delete(id *string) *model.Erro {
	if _, ok := (*db.ClientMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.ClientMap, *id)
	return nil
}

func (db MockClientRepository) Get(id *string) (interface{}, *model.Erro) {
	if clientResponse, ok := (*db.ClientMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &clientResponse, nil
	}

}

func (db MockClientRepository) Update(request interface{}) *model.Erro {
	clientRequest, ok := request.(*Client)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*db.ClientMap)[clientRequest.Client_id]; !ok {
		return model.IDnotFound
	}
	return nil
}

func (db MockClientRepository) GetAll() (interface{}, *model.Erro) {
	if len(*db.ClientMap) == 0 {
		return nil, model.IDnotFound
	}
	clients := make([]Client, 0, len(*db.ClientMap))
	for _, client := range *db.ClientMap {
		clients = append(clients, client)
	}
	return clients, nil
}
