package service

type AccountDB struct {
	account_id uint32
	client_id  uint32
	user_id    uint32
	agency_id  uint32
	password   string
	balance    float64
	status     bool
}

type TransferDB struct {
	account_id_from uint32
	account_id_to   uint32
	value           float64
}

type DepositDB struct {
	account_id uint32
	client_id  uint32
	user_id    uint32
	agency_iD  uint32
	deposit    float64
	balance    float64
}

type WithdrawalDB struct {
	account_id uint32
	client_id  uint32
	user_id    uint32
	agency_iD  uint32
	password   string
	withdrawal float64
	balance    float64
}

type ClientDB struct {
	client_id     uint32
	user_id       uint32
	name          string
	document      string
	password      string
	register_date string
	status        bool
}

type UserDB struct {
	user_id       uint32
	name          string
	document      string
	password      string
	register_date string
	status        bool
}
