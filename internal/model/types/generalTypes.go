package model

import "github.com/golang-jwt/jwt/v5"

type StandartResponse struct {
	Message string `json:"message" xml:"message"`
}

type Claims struct {
	Id   string `json:"id" xml:"id"`
	Role string `json:"role" xml:"role"`
	jwt.RegisteredClaims
}
