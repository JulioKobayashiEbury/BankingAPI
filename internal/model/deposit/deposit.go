package deposit

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo"
)

type DepositRepository interface {
	model.Repository[Deposit]
	GetFilteredByAccountID(context.Context, *string) (*[]Deposit, *echo.HTTPError)
}

type Deposit struct {
	Deposit_id   string  `json:"deposit_id" xml:"deposit_id"`
	Account_id   string  `json:"account_id" xml:"account_id"`
	Client_id    string  `json:"client_id" xml:"client_id"`
	User_id      string  `json:"user_id" xml:"user_id"`
	Agency_id    uint32  `json:"agency_id" xml:"agency_id"`
	Deposit      float64 `json:"deposit" xml:"deposit"`
	Deposit_date string  `json:"deposit_date" xml:"deposit_date"`
}
