package domain

type ClientRequest struct {
	ClientId     uint32 `json:"client_id" xml:"client_id"`
	UserId       uint32 `json:"user_id" xml:"user_id"`
	Name         string `json:"name" xml:"name"`
	Document     string `json:"document" xml:"document"`
	Password     string `json:"password" xml:"password"`
	RegisterDate string `json:"register_date" xml:"register_date"`
	Status       bool   `json:"status" xml:"status"`
}
