package model

import (
	"context"
)

const (
	AccountsPath    = "accounts"
	UsersPath       = "users"
	ClientPath      = "clients"
	TransfersPath   = "transfers"
	AutoDebit       = "autodebit"
	AutoDebitLog    = "autodebitlog"
	DepositPath     = "deposits"
	WithdrawalsPath = "withdrawals"
	CacheDuration   = 2
)

type Repository[T any] interface {
	Create(context.Context, *T) (*T, *Erro)
	Delete(context.Context, *string) *Erro
	Get(context.Context, *string) (*T, *Erro)
	Update(context.Context, *T) *Erro
	GetAll(context.Context) (*[]T, *Erro)
	GetFilteredByID(context.Context, *string) (*[]T, *Erro)
}

/* func TokenizeFilters(filters *string) *[]string {
	tokens := strings.Split(*filters, ",")
	return &tokens
}
*/
