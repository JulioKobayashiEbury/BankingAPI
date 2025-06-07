package automaticdebit

import "BankingAPI/internal/model"

var singleton *model.RepositoryInterface

type MockAutoDebitRepository struct {
	AutoDebitMap *map[string]AutomaticDebit
}

func NewMockAutoDebitRepository() model.RepositoryInterface {
	if singleton != nil {
		return *singleton
	}
	autoDebitMap := make(map[string]AutomaticDebit)
	*singleton = MockAutoDebitRepository{
		AutoDebitMap: &autoDebitMap,
	}
	return *singleton
}

func (db MockAutoDebitRepository) Create(request interface{}) (interface{}, *model.Erro) {
	autoDebitRequest, ok := request.(*AutomaticDebit)
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
	(*db.AutoDebitMap)[autoDebitRequest.Debit_id] = *autoDebitRequest
	return &autoDebitRequest, nil
}

func (db MockAutoDebitRepository) Delete(id *string) *model.Erro {
	if _, ok := (*db.AutoDebitMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.AutoDebitMap, *id)
	return nil
}

func (db MockAutoDebitRepository) Get(id *string) (interface{}, *model.Erro) {
	if autoDebit, ok := (*db.AutoDebitMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &autoDebit, nil
	}
}

func (db MockAutoDebitRepository) Update(request interface{}) *model.Erro {
	autoDebitRequest, ok := request.(*AutomaticDebit)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*db.AutoDebitMap)[autoDebitRequest.Debit_id]; !ok {
		return model.IDnotFound
	}
	(*db.AutoDebitMap)[autoDebitRequest.Debit_id] = *autoDebitRequest
	return nil
}
func (db MockAutoDebitRepository) GetAll() (interface{}, *model.Erro) {
	if len(*db.AutoDebitMap) == 0 {
		return nil, model.IDnotFound
	}
	autoDebits := make([]AutomaticDebit, 0, len(*db.AutoDebitMap))
	for _, autoDebit := range *db.AutoDebitMap {
		autoDebits = append(autoDebits, autoDebit)
	}
	return autoDebits, nil
}
