package account

import (
	"context"

	"BankingAPI/internal/model"

	"github.com/labstack/echo/v4"
)

type AccountRepository interface {
	model.Repository[Account]
	GetFilteredByClientID(context.Context, *string) (*[]Account, *echo.HTTPError)
}

type Account struct {
	Account_id    string       `json:"account_id" xml:"account_id"`
	Client_id     string       `json:"client_id" xml:"client_id"`
	User_id       string       `json:"user_id" xml:"user_id"`
	Agency_id     uint32       `json:"agency_id" xml:"agency_id"`
	Register_date string       `json:"register_date" xml:"register_date"`
	Balance       float64      `json:"balance" xml:"balance"`
	Status        model.Status `json:"status" xml:"status"`
}

type ListRequest struct {
	Filter string `json:"filter" xml:"filter"`
	Order  string `json:"order" xml:"order"`
}

type AccountReport struct {
	Account_id       string       `json:"account_id" xml:"account_id"`
	Client_id        string       `json:"client_id" xml:"client_id"`
	Agency_id        uint32       `json:"agency_id" xml:"agency_id"`
	Balance          float64      `json:"balance" xml:"balance"`
	Register_date    string       `json:"register_date" xml:"register_date"`
	Status           model.Status `json:"status" xml:"status"`
	Transfers        interface{}  `json:"transfer" xml:"transfers"`
	Deposits         interface{}  `json:"deposits" xml:"deposits"`
	Withdrawals      interface{}  `json:"withdrawals" xml:"withdrawals"`
	Automatic_Debits interface{}  `json:"automatic_debits" xml:"automatic_debits"`
	Report_Date      string       `json:"report_date" xml:"report_date"`
}
