package account

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo"
)

type MockAccountRepository struct {
	accountMap *map[string]Account
}

func NewMockAccountRepository() AccountRepository {
	return &MockAccountRepository{
		accountMap: &map[string]Account{},
	}
}

func (db MockAccountRepository) Create(ctx context.Context, request *Account) (*Account, *echo.HTTPError) {
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
	(*db.accountMap)[request.Account_id] = *request
	return request, nil
}

func (db MockAccountRepository) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if _, ok := (*db.accountMap)[*id]; !ok {
		return model.ErrIDnotFound
	}
	delete(*db.accountMap, *id)
	return nil
}

func (db MockAccountRepository) Get(ctx context.Context, id *string) (*Account, *echo.HTTPError) {
	if account, ok := (*db.accountMap)[*id]; !ok {
		return nil, model.ErrIDnotFound
	} else {
		return &account, nil
	}
}

func (db MockAccountRepository) Update(ctx context.Context, request *Account) *echo.HTTPError {
	if _, ok := (*db.accountMap)[request.Account_id]; !ok {
		return model.ErrIDnotFound
	}
	(*db.accountMap)[request.Account_id] = *request
	return nil
}

func (db MockAccountRepository) GetAll(ctx context.Context) (*[]Account, *echo.HTTPError) {
	if len(*db.accountMap) == 0 {
		return nil, model.ErrIDnotFound
	}
	accounts := make([]Account, 0, len(*db.accountMap))
	for _, account := range *db.accountMap {
		accounts = append(accounts, account)
	}
	return &accounts, nil
}

func (db MockAccountRepository) GetFilteredByClientID(ctx context.Context, clientID *string) (*[]Account, *echo.HTTPError) {
	if clientID == nil || len(*clientID) == 0 {
		return nil, model.ErrFilterNotSet
	}

	accountSlice := make([]Account, 0, len(*db.accountMap))

	for _, account := range *db.accountMap {
		if account.Client_id == *clientID {
			accountSlice = append(accountSlice, account)
		}
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
	return &accountSlice, nil
}
