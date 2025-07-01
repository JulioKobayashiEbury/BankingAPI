package account

import (
	"context"

	"BankingAPI/internal/model"
)

var singleton *AccountRepository

type MockAccountRepository struct {
	AccountMap *map[string]Account
}

func NewMockAccountRepository() *AccountRepository {
	if singleton != nil {
		return singleton
	}
	accountMap := make(map[string]Account)
	*singleton = MockAccountRepository{
		AccountMap: &accountMap,
	}
	return singleton
}

func (db MockAccountRepository) Create(ctx context.Context, request *Account) (*Account, *model.Erro) {
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
	(*db.AccountMap)[request.Account_id] = *request
	return request, nil
}

func (db MockAccountRepository) Delete(ctx context.Context, id *string) *model.Erro {
	if _, ok := (*db.AccountMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.AccountMap, *id)
	return nil
}

func (db MockAccountRepository) Get(ctx context.Context, id *string) (*Account, *model.Erro) {
	if account, ok := (*db.AccountMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &account, nil
	}
}

func (db MockAccountRepository) Update(ctx context.Context, request *Account) *model.Erro {
	if _, ok := (*db.AccountMap)[request.Account_id]; !ok {
		return model.IDnotFound
	}
	(*db.AccountMap)[request.Account_id] = *request
	return nil
}

func (db MockAccountRepository) GetAll(ctx context.Context) (*[]Account, *model.Erro) {
	if len(*db.AccountMap) == 0 {
		return nil, model.IDnotFound
	}
	accounts := make([]Account, 0, len(*db.AccountMap))
	for _, account := range *db.AccountMap {
		accounts = append(accounts, account)
	}
	return &accounts, nil
}

func (db MockAccountRepository) GetFilteredByID(ctx context.Context, filters *string) (*[]Account, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}

	/* accountSlice := make([]Account, 0, len(*db.AccountMap))
	for _, account := range *db.AccountMap {
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
			case "client_id":
				if operator == "==" && account.Client_id != value {
					match = false
				}
			case "status":
				if operator == "==" && account.Status != model.Status(value) {
					match = false
				}
			// Add more fields as necessary
			default:
				match = true
			}
		}
		if match {
			accountSlice = append(accountSlice, account)
		}
	}
	*/
	return nil, nil
}
