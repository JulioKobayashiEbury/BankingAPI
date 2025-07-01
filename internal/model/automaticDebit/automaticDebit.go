package automaticdebit

import "BankingAPI/internal/model"

type AutoDebitRepository interface {
	model.Repository[AutomaticDebit]
}
type AutomaticDebit struct {
	Debit_id        string  `json:"debit_id" xml:"debit_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	User_id         string  `json:"user_id" xml:"user_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Value           float64 `json:"value" xml:"value"`
	Debit_day       uint16  `json:"debit_day" xml:"debit_day"`
	Expiration_date string  `json:"expiration_date" xml:"expiration_date"`
	Register_date   string  `json:"register_date" xml:"register_date"`
}
