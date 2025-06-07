package transfer

import "BankingAPI/internal/model"

var singleton *model.RepositoryInterface

type MockTransferRepository struct {
	TransferMap *map[string]Transfer
}

func NewMockTransferRepository() model.RepositoryInterface {
	if singleton != nil {
		return *singleton
	}
	userMap := make(map[string]Transfer)
	*singleton = MockTransferRepository{
		TransferMap: &userMap,
	}
	return *singleton
}

func (t MockTransferRepository) Create(request interface{}) (interface{}, *model.Erro) {
	transferRequest, ok := request.(*Transfer)
	if !ok {
		return nil, model.DataTypeWrong
	}
	/*
		for {
			transferID := randomstring.String(10)
			if _, ok := (*m.TransferMap)[tranferID]; !ok {
				transferRequest.Transfer_id = transferID
				(*t.TransferMap)[transferRequest.Transfer_id] = *transferRequest
				break
			}
		}
	*/
	(*t.TransferMap)[transferRequest.Transfer_id] = *transferRequest
	return &transferRequest, nil
}

func (t MockTransferRepository) Delete(id *string) *model.Erro {
	if _, ok := (*t.TransferMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*t.TransferMap, *id)
	return nil
}

func (t MockTransferRepository) Get(id *string) (interface{}, *model.Erro) {
	if transfer, ok := (*t.TransferMap)[*id]; !ok {
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
	if _, ok := (*t.TransferMap)[transferRequest.Transfer_id]; !ok {
		return model.IDnotFound
	}
	(*t.TransferMap)[transferRequest.Transfer_id] = *transferRequest
	return nil
}

func (t MockTransferRepository) GetAll() (interface{}, *model.Erro) {
	if len(*t.TransferMap) == 0 {
		return nil, model.IDnotFound
	}
	users := make([]Transfer, 0, len(*t.TransferMap))
	for _, user := range *t.TransferMap {
		users = append(users, user)
	}
	return users, nil
}
