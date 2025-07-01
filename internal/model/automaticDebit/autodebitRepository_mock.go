package automaticdebit

import (
	"context"

	"BankingAPI/internal/model"
)

var singleton *AutoDebitRepository

type MockAutoDebitRepository struct {
	AutoDebitMap *map[string]AutomaticDebit
}

func NewMockAutoDebitRepository() AutoDebitRepository {
	if singleton != nil {
		return *singleton
	}
	autoDebitMap := make(map[string]AutomaticDebit)
	*singleton = MockAutoDebitRepository{
		AutoDebitMap: &autoDebitMap,
	}
	return *singleton
}

func (db MockAutoDebitRepository) Create(ctx context.Context, request *AutomaticDebit) (*AutomaticDebit, *model.Erro) {
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
	(*db.AutoDebitMap)[request.Debit_id] = *request
	return request, nil
}

func (db MockAutoDebitRepository) Delete(ctx context.Context, id *string) *model.Erro {
	if _, ok := (*db.AutoDebitMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.AutoDebitMap, *id)
	return nil
}

func (db MockAutoDebitRepository) Get(ctx context.Context, id *string) (*AutomaticDebit, *model.Erro) {
	if autoDebit, ok := (*db.AutoDebitMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &autoDebit, nil
	}
}

func (db MockAutoDebitRepository) Update(ctx context.Context, request *AutomaticDebit) *model.Erro {
	if _, ok := (*db.AutoDebitMap)[request.Debit_id]; !ok {
		return model.IDnotFound
	}
	(*db.AutoDebitMap)[request.Debit_id] = *request
	return nil
}

func (db MockAutoDebitRepository) GetAll(ctx context.Context) (*[]AutomaticDebit, *model.Erro) {
	if len(*db.AutoDebitMap) == 0 {
		return nil, model.IDnotFound
	}
	autoDebits := make([]AutomaticDebit, 0, len(*db.AutoDebitMap))
	for _, autoDebit := range *db.AutoDebitMap {
		autoDebits = append(autoDebits, autoDebit)
	}
	return &autoDebits, nil
}

func (db MockAutoDebitRepository) GetFilteredByID(ctx context.Context, filters *string) (*[]AutomaticDebit, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	autodebitSlice := make([]AutomaticDebit, 0, len(*db.AutoDebitMap))
	/* for _, autodebit := range *db.AutoDebitMap {
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
			default:
				match = false
			}
		}
		if match {
			autodebitSlice = append(autodebitSlice, autodebit)
		}
	}
	*/
	return &autodebitSlice, nil
}
