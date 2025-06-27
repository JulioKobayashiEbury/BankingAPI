package model

import (
	"strings"
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

type RepositoryInterface interface {
	Create(interface{}) (interface{}, *Erro)
	Delete(*string) *Erro
	Get(id *string) (interface{}, *Erro)
	Update(interface{}) *Erro
	GetAll() (interface{}, *Erro)
	GetFiltered(*[]string) (interface{}, *Erro)
}

func TokenizeFilters(filters *string) *[]string {
	tokens := strings.Split(*filters, ",")
	return &tokens
}
