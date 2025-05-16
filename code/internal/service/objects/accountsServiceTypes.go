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

type Transfer struct {
	account_id_from uint32
	account_id_to   uint32
	value           float64
}

type Deposit struct {
	account_id uint32
	client_id  uint32
	user_id    uint32
	agency_iD  uint32
	deposit    float64
	balance    float64
}

type Withdrawal struct {
	Account_id uint32
	Client_id  uint32
	User_id    uint32
	Agency_iD  uint32
	Password   string
	Withdrawal float64
	balance    float64
}
