package deposit

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo/v4"
)

var singleton *DepositRepository

type MockDepositRepository struct {
	DepositMap *map[string]Deposit
}

func NewMockDepositRepository() DepositRepository {
	if singleton != nil {
		return *singleton
	}
	userMap := make(map[string]Deposit)
	*singleton = MockDepositRepository{
		DepositMap: &userMap,
	}
	return *singleton
}

func (d MockDepositRepository) Create(ctx context.Context, request *Deposit) (*Deposit, *echo.HTTPError) {
	/*
		for {
			depositID := randomstring.String(10)
			if _, ok := (*d.DepositMap)[depositID]; !ok {
				depositRequest.Deposit_id = depositID
				(*t.DepositMap)[depositRequest.Deposit_id] = *depositRequest
				break
			}
		}
	*/
	(*d.DepositMap)[request.Deposit_id] = *request
	return request, nil
}

func (d MockDepositRepository) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if _, ok := (*d.DepositMap)[*id]; !ok {
		return model.ErrIDnotFound
	}
	delete(*d.DepositMap, *id)
	return nil
}

func (d MockDepositRepository) Get(ctx context.Context, id *string) (*Deposit, *echo.HTTPError) {
	if deposit, ok := (*d.DepositMap)[*id]; !ok {
		return nil, model.ErrIDnotFound
	} else {
		return &deposit, nil
	}
}

func (d MockDepositRepository) Update(ctx context.Context, request *Deposit) *echo.HTTPError {
	if _, ok := (*d.DepositMap)[request.Deposit_id]; !ok {
		return model.ErrIDnotFound
	}
	(*d.DepositMap)[request.Deposit_id] = *request
	return nil
}

func (d MockDepositRepository) GetAll(ctx context.Context) (*[]Deposit, *echo.HTTPError) {
	if len(*d.DepositMap) == 0 {
		return nil, model.ErrIDnotFound
	}
	users := make([]Deposit, 0, len(*d.DepositMap))
	for _, user := range *d.DepositMap {
		users = append(users, user)
	}
	return &users, nil
}

func (db MockDepositRepository) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]Deposit, *echo.HTTPError) {
	if accountID == nil || len(*accountID) == 0 {
		return nil, model.ErrFilterNotSet
	}
	depositSlice := make([]Deposit, 0, len(*db.DepositMap))

	for _, deposit := range *db.DepositMap {
		if deposit.Account_id == *accountID {
			depositSlice = append(depositSlice, deposit)
		}
	}
	/* for _, depositResponse := range *db.DepositMap {
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
				if operator == "==" && depositResponse.Account_id != value {
					match = false
				}
			case "client_id":
				if operator == "==" && depositResponse.Client_id != value {
					match = false
				}
			// Add more fields as necessary
			default:
				match = false
			}
		}
		if match {
			depositSlice = append(depositSlice, depositResponse)
		}
	}
	*/
	return &depositSlice, nil
}
