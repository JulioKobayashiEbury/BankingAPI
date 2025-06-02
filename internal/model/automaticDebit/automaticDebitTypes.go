package automaticdebit

type AutomaticDebit struct {
	Debit_id        string  `json:"debit_id" xml:"debit_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	Client_id       string  `json:"client_id" xml:"client_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Value           float64 `json:"value" xml:"value"`
	Debit_day       uint16  `json:"debit_day" xml:"debit_day"`
	Status          bool    `json:"status" xml:"status"`
	Expiration_date string  `json:"expiration_date" xml:"expiration_date"`
}
