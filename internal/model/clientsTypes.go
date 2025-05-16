package model

type ClientRequest struct {
	Client_id uint32 `json:"client_id" xml:"client_id"`
	User_id   uint32 `json:"user_id" xml:"user_id"`
	Name      string `json:"name" xml:"name"`
	Document  string `json:"document" xml:"document"`
	Password  string `json:"password" xml:"password"`
}

type ClientResponse struct {
	Client_id     uint32 `json:"client_id" xml:"client_id"`
	User_id       uint32 `json:"user_id" xml:"user_id"`
	Name          string `json:"name" xml:"name"`
	Document      string `json:"document" xml:"document"`
	Register_date string `json:"register_date" xml:"register_date"`
	Status        bool   `json:"status" xml:"status"`
}
