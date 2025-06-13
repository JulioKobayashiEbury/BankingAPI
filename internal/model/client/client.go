package client

type Client struct {
	Client_id     string `json:"client_id" xml:"client_id"`
	User_id       string `json:"user_id" xml:"user_id"`
	Name          string `json:"name" xml:"name"`
	Document      string `json:"document" xml:"document"`
	Register_date string `json:"register_date" xml:"register_date"`
	Status        bool   `json:"status" xml:"status"`
}

type ClientReport struct {
	Client_id     string      `json:"client_id" xml:"client_id"`
	User_id       string      `json:"user_id" xml:"user_id"`
	Name          string      `json:"name" xml:"name"`
	Document      string      `json:"document" xml:"document"`
	Register_date string      `json:"register_date" xml:"register_date"`
	Status        bool        `json:"status" xml:"status"`
	Accounts      interface{} `json:"accounts" xml:"accounts"`
	Report_date   string      `json:"report_date" xml:"report_date"`
}
