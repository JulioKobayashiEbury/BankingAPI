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
	UserDatabase           RepositoryInterface
	ClientDatabase         RepositoryInterface
	AccountDatabase        RepositoryInterface
	AutomaticDebitDatabase RepositoryInterface
	DepositDatabase        RepositoryInterface
	TransferDatabase       RepositoryInterface
	WithdrawalDatabase     RepositoryInterface
}

type StandartResponse struct {
	Message string `json:"message" xml:"message"`
}

type Claims struct {
	Id   string `json:"id" xml:"id"`
	Role string `json:"role" xml:"role"`
	jwt.RegisteredClaims
}

type Erro struct {
	Err      error
	HttpCode int
}

type Status string

func (s Status) IsValid() bool {
	for _, status := range ValidStatus {
		if s == status {
			return true
		}
	}
	return false
}
