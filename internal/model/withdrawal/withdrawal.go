package withdrawal

type Withdrawal struct {
	Withdrawal_id   string  `json:"withdrawal_id" xml:"withdrawal_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	Client_id       string  `json:"client_id" xml:"client_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Withdrawal      float64 `json:"withdrawal" xml:"withdrawal"`
	Withdrawal_date string  `json:"withdrawal_date" xml:"withdrawal_date"`
	Status          bool    `json:"status" xml:"status"`
}
