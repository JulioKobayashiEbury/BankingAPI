package transfer

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
