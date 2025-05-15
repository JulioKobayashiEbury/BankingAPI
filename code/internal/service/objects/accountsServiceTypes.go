package service

type Account struct {
	account_id uint32
	client_id  uint32
	user_id    uint32
	agency_id  uint32
	password   string
	balance    float64
	status     bool
}
