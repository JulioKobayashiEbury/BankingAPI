package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	//"2006-01-02T15:04:05+07:00"
	TimeLayout = time.RFC3339
)

type RepositoryList struct {
	UserDatabase    RepositoryInterface
	ClientDatabase  RepositoryInterface
	AccountDatabase RepositoryInterface
}

type StandartResponse struct {
	Message string `json:"message" xml:"message"`
}

type Claims struct {
	Id   string `json:"id" xml:"id"`
	Role string `json:"role" xml:"role"`
	jwt.RegisteredClaims
}
