package controller

type AccountRequest struct {
	Account_id uint32 `json:"account_id" xml:"account_id"`
	Client_id  uint32 `json:"client_id" xml:"client_id"`
	User_id    uint32 `json:"user_id" xml:"user_id"`
	Agency_id  uint32 `json:"agency_id" xml:"agency_id"`
	Password   string `json:"password" xml:"password"`
}

type AccountResponse struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Client_id  uint32  `json:"client_id" xml:"client_id"`
	User_id    uint32  `json:"user_id" xml:"user_id"`
	Agency_id  uint32  `json:"agency_id" xml:"agency_id"`
	Balance    float64 `json:"balance" xml:"balance"`
	Status     bool    `json:"status" xml:"status"`
}

type WithdrawalRequest struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Client_id  uint32  `json:"client_id" xml:"client_id"`
	User_id    uint32  `json:"user_id" xml:"user_id"`
	Agency_iD  uint32  `json:"agency_id" xml:"agency_id"`
	Password   string  `json:"password" xml:"password"`
	Withdrawal float64 `json:"withdrawal" xml:"withdrawal"`
}

type DepositRequest struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Client_id  uint32  `json:"client_id" xml:"client_id"`
	User_id    uint32  `json:"user_id" xml:"user_id"`
	Agency_iD  uint32  `json:"agency_id" xml:"agency_id"`
	Deposit    float64 `json:"deposit" xml:"deposit"`
}

type DepositResponse struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Balance    float64 `json:"balance" xml:"balance"`
}

type WithdrawalResponse struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Balance    float64 `json:"balance" xml:"balance"`
}
type TransferRequest struct {
	Account_id_from uint32  `json:"account_id_from" xml:"account_id_from"`
	Account_id_to   uint32  `json:"account_id_to" xml:"account_id_to"`
	Password        string  `json:"password" xml:"password"`
	Value           float64 `json:"value" xml:"value"`
	User_id_from    uint32  `json:"user_id_from" xml:"user_id_from"`
	User_id_to      uint32  `json:"user_id_to" xml:"user_id_to"`
}

type AutomaticDebitRequest struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Client_id  uint32  `json:"client_id" xml:"client_id"`
	Agency_iD  uint32  `json:"agency_id" xml:"agency_id"`
	Password   string  `json:"password" xml:"password"`
	Value      float64 `json:"value" xml:"value"`
	Debit_date string  `json:"dabit_date" xml:"debit_date"`
}

type ListRequest struct {
	Filter string `json:"filter" xml:"filter"`
	Order  string `json:"order" xml:"order"`
}

type StandartResponse struct {
	Message string `json:"message" xml:"message"`
}
