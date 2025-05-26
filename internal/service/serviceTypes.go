package service

const (
	UserRole    = "user"
	ClientRole  = "client"
	AccountRole = "account"
)

type AccountDB struct {
	account_id    string
	client_id     string
	user_id       string
	agency_id     uint32
	password      string
	register_date string
	balance       float64
	status        bool
}

type TransferDB struct {
	transfer_id   string
	account_id    string
	account_to    string
	value         float64
	transfer_date string
}

type DepositDB struct {
	deposit_id   string
	account_id   string
	client_id    string
	user_id      string
	agency_id    uint32
	deposit      float64
	deposit_date string
}

type WithdrawalDB struct {
	withdrawal_id   string
	account_id      string
	client_id       string
	user_id         string
	agency_id       uint32
	withdrawal      float64
	withdrawal_date string
}

type ClientDB struct {
	client_id     string
	user_id       string
	name          string
	document      string
	password      string
	register_date string
	status        bool
}

type UserDB struct {
	user_id       string
	name          string
	document      string
	password      string
	register_date string
	status        bool
}

type Auth struct {
	Password string `json:"password" xml:"passsword"`
}
