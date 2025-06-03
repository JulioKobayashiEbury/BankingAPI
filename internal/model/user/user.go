package user

import "BankingAPI/internal/model/client"

type User struct {
	User_id       string `json:"user_id" xml:"user_id"`
	Name          string `json:"name" xml:"name" validate:"required"`
	Document      string `json:"document" xml:"document" validate:"required"`
	Password      string `json:"password" xml:"password" validate:"required"`
	Register_date string `json:"register_date" xml:"register_date"`
	Status        bool   `json:"status" xml:"status"`
}

type UserRequest struct {
	User_id  string `json:"user_id" xml:"user_id"`
	Name     string `json:"name" xml:"name" validate:"required"`
	Document string `json:"document" xml:"document" validate:"required"`
	Password string `json:"password" xml:"password" validate:"required"`
}

type UserResponse struct {
	User_id       string `json:"user_id" xml:"user_id"`
	Name          string `json:"name" xml:"name" validate:"required"`
	Document      string `json:"document" xml:"document" validate:"required"`
	Register_date string `json:"register_date" xml:"register_date"`
	Status        bool   `json:"status" xml:"status"`
}

type UserReport struct {
	User_id       string                  `json:"user_id" xml:"user_id"`
	Name          string                  `json:"name" xml:"name" validate:"required"`
	Document      string                  `json:"document" xml:"document" validate:"required"`
	Register_date string                  `json:"register_date" xml:"register_date"`
	Status        bool                    `json:"status" xml:"status"`
	Clients       []client.ClientResponse `json:"clients" xml:"clients"`
	Report_date   string                  `json:"report_date" xml:"report_date"`
}

type AuthenticateUser struct {
	User_id  string `json:"user_id" xml:"user_id"`
	Password string `json:"password" xml:"password" validate:"required"`
}
