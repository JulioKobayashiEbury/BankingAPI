package deposit

import "BankingAPI/internal/model"

var singleton *model.RepositoryInterface

type MockDepositRepository struct {
	DepositMap *map[string]Deposit
}

func NewMockDepositRepository() model.RepositoryInterface {
	if singleton != nil {
		return *singleton
	}
	userMap := make(map[string]Deposit)
	*singleton = MockDepositRepository{
		DepositMap: &userMap,
	}
	return *singleton
}

func (d MockDepositRepository) Create(request interface{}) (interface{}, *model.Erro) {
	depositRequest, ok := request.(*Deposit)
	if !ok {
		return nil, model.DataTypeWrong
	}
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
	(*d.DepositMap)[depositRequest.Deposit_id] = *depositRequest
	return &depositRequest, nil

}

func (d MockDepositRepository) Delete(id *string) *model.Erro {
	if _, ok := (*d.DepositMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*d.DepositMap, *id)
	return nil
}

func (d MockDepositRepository) Get(id *string) (interface{}, *model.Erro) {
	if deposit, ok := (*d.DepositMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &deposit, nil
	}
}

func (d MockDepositRepository) Update(request interface{}) *model.Erro {
	depositRequest, ok := request.(*Deposit)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*d.DepositMap)[depositRequest.Deposit_id]; !ok {
		return model.IDnotFound
	}
	(*d.DepositMap)[depositRequest.Deposit_id] = *depositRequest
	return nil
}
func (d MockDepositRepository) GetAll() (interface{}, *model.Erro) {
	if len(*d.DepositMap) == 0 {
		return nil, model.IDnotFound
	}
	users := make([]Deposit, 0, len(*d.DepositMap))
	for _, user := range *d.DepositMap {
		users = append(users, user)
	}
	return users, nil
}
