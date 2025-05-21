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
	account_id_from string
	account_id_to   string
	value           float64
	password        string
}

type DepositDB struct {
	account_id string
	client_id  string
	user_id    string
	agency_id  uint32
	deposit    float64
	balance    float64
}

type WithdrawalDB struct {
	account_id string
	client_id  string
	user_id    string
	agency_id  uint32
	password   string
	withdrawal float64
	balance    float64
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
