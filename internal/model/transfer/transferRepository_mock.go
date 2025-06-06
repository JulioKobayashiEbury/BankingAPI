package transfer

import "BankingAPI/internal/model"

var singleton *model.RepositoryInterface

type MockTransferRepository struct {
	UserMap *map[string]Transfer
}

func NewMockTransferRepository() model.RepositoryInterface {
	if singleton != nil {
		return *singleton
	}
	userMap := make(map[string]Transfer)
	*singleton = MockTransferRepository{
		UserMap: &userMap,
	}
	return *singleton
}

func (t MockTransferRepository) Create(request interface{}) (*string, *model.Erro) {
	transferRequest, ok := request.(*Transfer)
	if !ok {
		return nil, model.DataTypeWrong
	}
	/*
		for {
			transferID := randomstring.String(10)
			if _, ok := (*m.UserMap)[tranferID]; !ok {
				transferRequest.Transfer_id = transferID
				(*t.UserMap)[transferRequest.Transfer_id] = *transferRequest
				break
			}
		}
	*/
	(*t.UserMap)[transferRequest.Transfer_id] = *transferRequest
	return &transferRequest.Transfer_id, nil
}

func (t MockTransferRepository) Delete(id *string) *model.Erro {
	if _, ok := (*t.UserMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*t.UserMap, *id)
	return nil
}

func (t MockTransferRepository) Get(id *string) (interface{}, *model.Erro) {
	if transfer, ok := (*t.UserMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &transfer, nil
	}
}

func (t MockTransferRepository) Update(request interface{}) *model.Erro {
	transferRequest, ok := request.(*Transfer)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*t.UserMap)[transferRequest.Transfer_id]; !ok {
		return model.IDnotFound
	}
	(*t.UserMap)[transferRequest.Transfer_id] = *transferRequest
	return nil
}

func (t MockTransferRepository) GetAll() (interface{}, *model.Erro) {
	if len(*t.UserMap) == 0 {
		return nil, model.IDnotFound
	}
	users := make([]Transfer, 0, len(*t.UserMap))
	for _, user := range *t.UserMap {
		users = append(users, user)
	}
	return users, nil
}
