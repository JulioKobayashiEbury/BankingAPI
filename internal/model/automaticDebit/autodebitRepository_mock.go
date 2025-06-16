package automaticdebit

import (
	"BankingAPI/internal/model"
)

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

func (db MockAutoDebitRepository) GetFiltered(filters *[]string) (interface{}, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	autodebitSlice := make([]AutomaticDebit, 0, len(*db.AutoDebitMap))
	for _, autodebit := range *db.AutoDebitMap {
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
			case "account_id":
				if operator == "==" && autodebit.Account_id != value {
					match = false
				}
			case "status":
				if operator == "==" && autodebit.Status != model.Status(value) {
					match = false
				}
			// Add more fields as necessary
			default:
				match = false
			}
		}
		if match {
			autodebitSlice = append(autodebitSlice, autodebit)
		}
	}
	return &autodebitSlice, nil
}
