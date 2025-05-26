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

type WithdrawalResponse struct {
	Withdrawal_id   string  `json:"withdrawal_id" xml:"withdrawal_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	Client_id       string  `json:"client_id" xml:"client_id"`
	Agency_iD       uint32  `json:"agency_id" xml:"agency_id"`
	Withdrawal      float64 `json:"withdrawal" xml:"withdrawal"`
	Withdrawal_date string  `json:"trasnfer_date" xml:"trasnfer_date"`
}

type DepositRequest struct {
	Account_id string  `json:"account_id" xml:"account_id"`
	Client_id  string  `json:"client_id" xml:"client_id"`
	User_id    string  `json:"user_id" xml:"user_id"`
	Agency_id  uint32  `json:"agency_id" xml:"agency_id"`
	Deposit    float64 `json:"deposit" xml:"deposit"`
}

type DepositResponse struct {
	Deposit_id   string  `json:"deposit_id" xml"deposit_id"`
	Account_id   string  `json:"account_id" xml:"account_id"`
	Client_id    string  `json:"client_id" xml:"client_id"`
	User_id      string  `json:"user_id" xml:"user_id"`
	Agency_id    uint32  `json:"agency_id" xml:"agency_id"`
	Deposit      float64 `json:"deposit" xml:"deposit"`
	Deposit_date string  `json:"deposit_date" xml:"daposit_date"`
}

type TransferRequest struct {
	Transfer_id string  `json:"transfer_id" xml:"transfer_id"`
	Account_id  string  `json:"account_id" xml:"account_id"`
	Account_to  string  `json:"account_to" xml:"account_to"`
	Value       float64 `json:"value" xml:"value"`
}

type TransferResponse struct {
	Transfer_id   string  `json:"transfer_id" xml:"transfer_id"`
	Account_id    string  `json:"account_id" xml:"account_id"`
	Account_to    string  `json:"account_to" xml:"account_to"`
	Value         float64 `json:"value" xml:"value"`
	Register_date string  `json:"register_date" xml:"register_date"`
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

type AccountReport struct {
	Account_id       string                   `json:"account_id" xml:"account_id"`
	Client_id        string                   `json:"client_id" xml:"client_id"`
	Agency_id        uint32                   `json:"agency_id" xml:"agency_id"`
	Balance          float64                  `json:"balance" xml:"balance"`
	Register_date    string                   `json:"register_date" xml:"register_date"`
	Status           bool                     `json:"status" xml:"status"`
	Transfers        []TransferResponse       `json:"transfer" xml:"transfers"`
	Deposits         []DepositResponse        `json:"deposits" xml:"deposits"`
	Withdrawals      []WithdrawalResponse     `json:"withdrawals" xml:"withdrawals"`
	Automatic_Debits []AutomaticDebitResponse `json:"automatic_debits" xml:"automatic_debits"`
	Report_Date      string                   `json:"report_date" xml:"report_date"`
}
