package deposit

import (
	"context"

	"BankingAPI/internal/model"
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

func (d MockDepositRepository) Create(ctx context.Context, request *Deposit) (*Deposit, *model.Erro) {
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

func (d MockDepositRepository) Delete(ctx context.Context, id *string) *model.Erro {
	if _, ok := (*d.DepositMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*d.DepositMap, *id)
	return nil
}

func (d MockDepositRepository) Get(ctx context.Context, id *string) (*Deposit, *model.Erro) {
	if deposit, ok := (*d.DepositMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &deposit, nil
	}
}

func (d MockDepositRepository) Update(ctx context.Context, request *Deposit) *model.Erro {
	if _, ok := (*d.DepositMap)[request.Deposit_id]; !ok {
		return model.IDnotFound
	}
	(*d.DepositMap)[request.Deposit_id] = *request
	return nil
}

func (d MockDepositRepository) GetAll(ctx context.Context) (*[]Deposit, *model.Erro) {
	if len(*d.DepositMap) == 0 {
		return nil, model.IDnotFound
	}
	users := make([]Deposit, 0, len(*d.DepositMap))
	for _, user := range *d.DepositMap {
		users = append(users, user)
	}
	return &users, nil
}

func (db MockDepositRepository) GetFilteredByID(ctx context.Context, filters *string) (*[]Deposit, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	depositSlice := make([]Deposit, 0, len(*db.DepositMap))
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
