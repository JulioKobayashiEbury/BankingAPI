package model

type AccountRequest struct {
	Account_id string `json:"account_id" xml:"account_id"`
	Client_id  string `json:"client_id" xml:"client_id"`
	User_id    string `json:"user_id" xml:"user_id"`
	Agency_id  uint32 `json:"agency_id" xml:"agency_id"`
	Password   string `json:"password" xml:"password"`
}

type AccountResponse struct {
	Account_id    string  `json:"account_id" xml:"account_id"`
	Client_id     string  `json:"client_id" xml:"client_id"`
	User_id       string  `json:"user_id" xml:"user_id"`
	Agency_id     uint32  `json:"agency_id" xml:"agency_id"`
	Register_date string  `json:"register_date" xml:"register_date"`
	Balance       float64 `json:"balance" xml:"balance"`
	Status        bool    `json:"status" xml:"status"`
}

type WithdrawalRequest struct {
	Account_id string  `json:"account_id" xml:"account_id"`
	Client_id  string  `json:"client_id" xml:"client_id"`
	Agency_iD  uint32  `json:"agency_id" xml:"agency_id"`
	Withdrawal float64 `json:"withdrawal" xml:"withdrawal"`
}

type DepositRequest struct {
	Account_id string  `json:"account_id" xml:"account_id"`
	Client_id  string  `json:"client_id" xml:"client_id"`
	User_id    string  `json:"user_id" xml:"user_id"`
	Agency_iD  uint32  `json:"agency_id" xml:"agency_id"`
	Deposit    float64 `json:"deposit" xml:"deposit"`
}

type DepositResponse struct {
	Account_id string  `json:"account_id" xml:"account_id"`
	Balance    float64 `json:"balance" xml:"balance"`
}

type WithdrawalResponse struct {
	Account_id string  `json:"account_id" xml:"account_id"`
	Balance    float64 `json:"balance" xml:"balance"`
}
type TransferRequest struct {
	Transfer_id string  `json:"transfer_id" xml:"transfer_id"`
	Account_id  string  `json:"account_id" xml:"account_id"`
	Account_to  string  `json:"account_to" xml:"account_to"`
	Password    string  `json:"password" xml:"password"`
	Value       float64 `json:"value" xml:"value"`
}

type AutomaticDebitRequest struct {
	Account_id      string  `json:"account_id" xml:"account_id"`
	Client_id       string  `json:"client_id" xml:"client_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Value           float64 `json:"value" xml:"value"`
	Debit_day       uint16  `json:"debit_day" xml:"debit_day"`
	Expiration_date string  `json:"expiration_date" xml:"expiration_date"`
}

type AutomaticDebitResponse struct {
	Debit_id        string  `json:"debit_id" xml:"debit_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	Client_id       string  `json:"client_id" xml:"client_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Value           float64 `json:"value" xml:"value"`
	Debit_day       uint16  `json:"debit_day" xml:"debit_day"`
	Expiration_date string  `json:"expiration_date" xml:"expiration_date"`
	Register_date   string  `json:"register_date" xml:"register_date"`
}

type ListRequest struct {
	Filter string `json:"filter" xml:"filter"`
	Order  string `json:"order" xml:"order"`
}
