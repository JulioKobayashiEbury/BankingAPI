package withdrawal

import "BankingAPI/internal/model"

type MockWithdrawalRepository struct {
	WithdrawalMap *map[string]Withdrawal
}

func (db MockWithdrawalRepository) Create(request interface{}) (interface{}, *model.Erro) {
	withdrawal, ok := request.(*Withdrawal)
	if !ok {
		return nil, model.DataTypeWrong
	}
	(*db.WithdrawalMap)[withdrawal.Withdrawal_id] = *withdrawal
	return db.Get(&withdrawal.Withdrawal_id)
}

func (db MockWithdrawalRepository) Delete(id *string) *model.Erro {
	if _, ok := (*db.WithdrawalMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*db.WithdrawalMap, *id)
	return nil
}

func (db MockWithdrawalRepository) Get(id *string) (interface{}, *model.Erro) {
	if withdrawal, ok := (*db.WithdrawalMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &withdrawal, nil
	}
}

func (db MockWithdrawalRepository) Update(request interface{}) *model.Erro {
	withdrawal, ok := request.(*Withdrawal)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*db.WithdrawalMap)[withdrawal.Withdrawal_id]; !ok {
		return model.IDnotFound
	}
	(*db.WithdrawalMap)[withdrawal.Withdrawal_id] = *withdrawal
	return nil
}

func (db MockWithdrawalRepository) GetAll() (interface{}, *model.Erro) {
	withdrawals := make([]Withdrawal, 0, len(*db.WithdrawalMap))
	for _, withdrawal := range *db.WithdrawalMap {
		withdrawals = append(withdrawals, withdrawal)
	}
	return &withdrawals, nil
}
