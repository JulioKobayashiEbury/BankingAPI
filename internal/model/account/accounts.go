package account

import (
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/withdrawal"
)

type AccountRequest struct {
	Account_id string `json:"account_id" xml:"account_id"`
	Client_id  string `json:"client_id" xml:"client_id"`
	User_id    string `json:"user_id" xml:"user_id"`
	Agency_id  uint32 `json:"agency_id" xml:"agency_id"`
}

type AccountResponse struct {
	Account_id       string   `json:"account_id" xml:"account_id"`
	Client_id        string   `json:"client_id" xml:"client_id"`
	User_id          string   `json:"user_id" xml:"user_id"`
	Agency_id        uint32   `json:"agency_id" xml:"agency_id"`
	Register_date    string   `json:"register_date" xml:"register_date"`
	Balance          float64  `json:"balance" xml:"balance"`
	Transfers        []string `json:"transfers" xml:"transfers"`
	Withdrawals      []string `json:"withdrawals" xml:"withdrawals"`
	Automatic_debits []string `json:"automatic_debits" xml:"automatic_debits"`
	Status           bool     `json:"status" xml:"status"`
}

type ListRequest struct {
	Filter string `json:"filter" xml:"filter"`
	Order  string `json:"order" xml:"order"`
}

type AccountReport struct {
	Account_id       string                          `json:"account_id" xml:"account_id"`
	Client_id        string                          `json:"client_id" xml:"client_id"`
	Agency_id        uint32                          `json:"agency_id" xml:"agency_id"`
	Balance          float64                         `json:"balance" xml:"balance"`
	Register_date    string                          `json:"register_date" xml:"register_date"`
	Status           bool                            `json:"status" xml:"status"`
	Transfers        []transfer.TransferResponse     `json:"transfer" xml:"transfers"`
	Deposits         []deposit.DepositResponse       `json:"deposits" xml:"deposits"`
	Withdrawals      []withdrawal.WithdrawalResponse `json:"withdrawals" xml:"withdrawals"`
	Automatic_Debits []automaticdebit.AutomaticDebit `json:"automatic_debits" xml:"automatic_debits"`
	Report_Date      string                          `json:"report_date" xml:"report_date"`
}
