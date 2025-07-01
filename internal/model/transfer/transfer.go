package transfer

import "BankingAPI/internal/model"

type TransferRepository interface {
	model.Repository[Transfer]
}

type Transfer struct {
	Transfer_id   string  `json:"transfer_id" xml:"transfer_id"`
	User_id       string  `json:"user_id" xml:"user_id"`
	Account_id    string  `json:"account_id" xml:"account_id"`
	User_to       string  `json:"user_to" xml:"user_to"`
	Account_to    string  `json:"account_to" xml:"account_to"`
	Value         float64 `json:"value" xml:"value"`
	Register_date string  `json:"register_date" xml:"register_date"`
}
