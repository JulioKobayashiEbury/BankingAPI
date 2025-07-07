package withdrawal

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo"
)

type WithdrawalRepository interface {
	model.Repository[Withdrawal]
	GetFilteredByAccountID(context.Context, *string) (*[]Withdrawal, *echo.HTTPError)
}
type Withdrawal struct {
	Withdrawal_id   string  `json:"withdrawal_id" xml:"withdrawal_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	User_id         string  `json:"user_id" xml:"user_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Withdrawal      float64 `json:"withdrawal" xml:"withdrawal"`
	Withdrawal_date string  `json:"withdrawal_date" xml:"withdrawal_date"`
}
