package automaticdebit

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo"
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

func (db MockAutoDebitRepository) Create(ctx context.Context, request *AutomaticDebit) (*AutomaticDebit, *echo.HTTPError) {
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

func (db MockAutoDebitRepository) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if _, ok := (*db.AutoDebitMap)[*id]; !ok {
		return model.ErrIDnotFound
	}
	delete(*db.AutoDebitMap, *id)
	return nil
}

func (db MockAutoDebitRepository) Get(ctx context.Context, id *string) (*AutomaticDebit, *echo.HTTPError) {
	if autoDebit, ok := (*db.AutoDebitMap)[*id]; !ok {
		return nil, model.ErrIDnotFound
	} else {
		return &autoDebit, nil
	}
}

func (db MockAutoDebitRepository) Update(ctx context.Context, request *AutomaticDebit) *echo.HTTPError {
	if _, ok := (*db.AutoDebitMap)[request.Debit_id]; !ok {
		return model.ErrIDnotFound
	}
	(*db.AutoDebitMap)[request.Debit_id] = *request
	return nil
}

func (db MockAutoDebitRepository) GetAll(ctx context.Context) (*[]AutomaticDebit, *echo.HTTPError) {
	if len(*db.AutoDebitMap) == 0 {
		return nil, model.ErrIDnotFound
	}
	autoDebits := make([]AutomaticDebit, 0, len(*db.AutoDebitMap))
	for _, autoDebit := range *db.AutoDebitMap {
		autoDebits = append(autoDebits, autoDebit)
	}
	return &autoDebits, nil
}

func (db MockAutoDebitRepository) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]AutomaticDebit, *echo.HTTPError) {
	if accountID == nil || len(*accountID) == 0 {
		return nil, model.ErrFilterNotSet
	}
	autodebitSlice := make([]AutomaticDebit, 0, len(*db.AutoDebitMap))

	for _, autodebit := range *db.AutoDebitMap {
		if autodebit.Account_id == *accountID {
			autodebitSlice = append(autodebitSlice, autodebit)
		}
	}
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
