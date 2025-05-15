package controller

type AccountRequest struct {
	Account_id uint32  `json:"account_id" xml:"account_id"`
	Client_id  uint32  `json:"client_id" xml:"client_id"`
	User_id    uint32  `json:"user_id" xml:"user_id"`
	Agency_iD  uint32  `json:"agency_id" xml:"agency_id"`
	Password   string  `json:"password" xml:"password"`
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
