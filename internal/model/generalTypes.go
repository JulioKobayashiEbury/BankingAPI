package model

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	//"2006-01-02T15:04:05+07:00"
	TimeLayout = time.RFC3339
)

var (
	InvalidFilterFormat = &Erro{Err: errors.New("repository Error: Invalid fitler format"), HttpCode: http.StatusBadRequest}
	FilterNotSet        = &Erro{Err: errors.New("repository Error: filter value not set"), HttpCode: http.StatusBadRequest}
	ResquestNotSet      = &Erro{Err: errors.New("repository Error: Request value not set"), HttpCode: http.StatusBadRequest}
	FailCreatingClient  = &Erro{Err: errors.New("repository Error: Failed to create DB client"), HttpCode: http.StatusInternalServerError}
	IDnotFound          = &Erro{Err: errors.New("repository Error: Id not found"), HttpCode: http.StatusBadRequest}
	DataTypeWrong       = &Erro{Err: errors.New("repository Error: Invalid argument passed"), HttpCode: http.StatusBadRequest}
	InvalidStatus       = &Erro{Err: errors.New("invalid status value"), HttpCode: http.StatusBadRequest}

	ValidStatus = []Status{"active", "blocked"}
)

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
