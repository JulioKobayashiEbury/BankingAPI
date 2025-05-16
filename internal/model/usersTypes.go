package model

type UserRequest struct {
	User_id  uint32 `json:"user_id" xml:"user_id"`
	Name     string `json:"name" xml:"name" validate:"required"`
	Document string `json:"document" xml:"document" validate:"required"`
	Password string `json:"password" xml:"password" validate:"required"`
}

type UserResponse struct {
	User_id       uint32 `json:"user_id" xml:"user_id"`
	Name          string `json:"name" xml:"name" validate:"required"`
	Document      string `json:"document" xml:"document" validate:"required"`
	Register_date string `json:"register_date" xml:"register_date"`
	Status        bool   `json:"status" xml:"status"`
}
