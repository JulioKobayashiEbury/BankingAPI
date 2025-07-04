package withdrawal

import "BankingAPI/internal/model"

type WithdrawalRepository interface {
	model.Repository[Withdrawal]
}
type Withdrawal struct {
	Withdrawal_id   string  `json:"withdrawal_id" xml:"withdrawal_id"`
	Account_id      string  `json:"account_id" xml:"account_id"`
	User_id         string  `json:"user_id" xml:"user_id"`
	Agency_id       uint32  `json:"agency_id" xml:"agency_id"`
	Withdrawal      float64 `json:"withdrawal" xml:"withdrawal"`
	Withdrawal_date string  `json:"withdrawal_date" xml:"withdrawal_date"`
}
