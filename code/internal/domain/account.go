package domain

type AccountRequest struct {
	AccountId       uint32  `json:"account_id" xml:"account_id"`
	ClientId        uint32  `json:"client_id" xml:"client_id" validate:"required"`
	UserId          uint32  `json:"user_id" xml:"user_id" validate:"required"`
	AgencyID        uint32  `json:"agency_id" xml:"agency_id" validate:"required"`
	AccountPassword string  `json:"password" xml:"password" validate:"required"`
	AccountBalance  float64 `json:"balance" xml:"balance"`
	Status          bool    `json:"status" xml:"status"`
}
