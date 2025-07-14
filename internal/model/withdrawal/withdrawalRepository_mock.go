package withdrawal

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo/v4"
)

type MockWithdrawalRepository struct {
	WithdrawalMap *map[string]Withdrawal
}

func NewMockWithdrawalRepository() WithdrawalRepository {
	return &MockWithdrawalRepository{
		WithdrawalMap: &map[string]Withdrawal{},
	}
}

func (db MockWithdrawalRepository) Create(ctx context.Context, request *Withdrawal) (*Withdrawal, *echo.HTTPError) {
	(*db.WithdrawalMap)[request.Withdrawal_id] = *request
	return db.Get(ctx, &request.Withdrawal_id)
}

func (db MockWithdrawalRepository) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if _, ok := (*db.WithdrawalMap)[*id]; !ok {
		return model.ErrIDnotFound
	}
	delete(*db.WithdrawalMap, *id)
	return nil
}

func (db MockWithdrawalRepository) Get(ctx context.Context, id *string) (*Withdrawal, *echo.HTTPError) {
	if withdrawal, ok := (*db.WithdrawalMap)[*id]; !ok {
		return nil, model.ErrIDnotFound
	} else {
		return &withdrawal, nil
	}
}

func (db MockWithdrawalRepository) Update(ctx context.Context, request *Withdrawal) *echo.HTTPError {
	if _, ok := (*db.WithdrawalMap)[request.Withdrawal_id]; !ok {
		return model.ErrIDnotFound
	}
	(*db.WithdrawalMap)[request.Withdrawal_id] = *request
	return nil
}

func (db MockWithdrawalRepository) GetAll(ctx context.Context) (*[]Withdrawal, *echo.HTTPError) {
	withdrawals := make([]Withdrawal, 0, len(*db.WithdrawalMap))
	for _, withdrawal := range *db.WithdrawalMap {
		withdrawals = append(withdrawals, withdrawal)
	}
	return &withdrawals, nil
}

func (db MockWithdrawalRepository) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]Withdrawal, *echo.HTTPError) {
	if accountID == nil || len(*accountID) == 0 {
		return nil, model.ErrFilterNotSet
	}

	withdrawalSlice := make([]Withdrawal, 0, len(*db.WithdrawalMap))
	for _, withdrawal := range *db.WithdrawalMap {
		if withdrawal.Account_id == *accountID {
			withdrawalSlice = append(withdrawalSlice, withdrawal)
		}
	}
	return &withdrawalSlice, nil
}
