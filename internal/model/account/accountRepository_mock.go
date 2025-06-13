package account

import (
	"BankingAPI/internal/model"
	"strconv"
)

var singleton *model.RepositoryInterface

type MockAccountRepository struct {
	AccountMap *map[string]Account
}

func NewMockAccountRepository() *model.RepositoryInterface {
	if singleton != nil {
		return singleton
	}
	accountMap := make(map[string]Account)
	*singleton = MockAccountRepository{
		AccountMap: &accountMap,
	}
	return singleton
}

func (db MockAccountRepository) Create(request interface{}) (interface{}, *model.Erro) {
	accountRequest, ok := request.(*Account)
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
	(*db.AccountMap)[accountRequest.Account_id] = *accountRequest
	return &accountRequest, nil
}

func (db MockAccountRepository) Delete(id *string) *model.Erro {
	if _, ok := (*db.AccountMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.AccountMap, *id)
	return nil
}

func (db MockAccountRepository) Get(id *string) (interface{}, *model.Erro) {
	if account, ok := (*db.AccountMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &account, nil
	}
}

func (db MockAccountRepository) Update(request interface{}) *model.Erro {
	accountRequest, ok := request.(*Account)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*db.AccountMap)[accountRequest.Account_id]; !ok {
		return model.IDnotFound
	}
	(*db.AccountMap)[accountRequest.Account_id] = *accountRequest
	return nil
}
func (db MockAccountRepository) GetAll() (interface{}, *model.Erro) {
	if len(*db.AccountMap) == 0 {
		return nil, model.IDnotFound
	}
	accounts := make([]Account, 0, len(*db.AccountMap))
	for _, account := range *db.AccountMap {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (db MockAccountRepository) GetFiltered(filters *[]string) (interface{}, *model.Erro) {

	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}

	accountSlice := make([]Account, 0, len(*db.AccountMap))
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
				floatValue, err := strconv.ParseBool(value)
				if err != nil {
					return nil, model.DataTypeWrong
				}
				if operator == "==" && account.Status != floatValue {
					match = false
				}
			// Add more fields as necessary
			default:
				match = false
			}
		}
		if match {
			accountSlice = append(accountSlice, account)
		}
	}

	return nil, nil
}
