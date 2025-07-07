package model

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
)

const (
	//"2006-01-02T15:04:05+07:00"
	TimeLayout = time.RFC3339
)

var (
	ErrInvalidFilterFormat = &echo.HTTPError{Internal: errors.New("repository Error: Invalid fitler format"), Code: http.StatusBadRequest, Message: "repository Error: Invalid fitler format"}
	ErrFilterNotSet        = &echo.HTTPError{Internal: errors.New("repository Error: filter value not set"), Code: http.StatusBadRequest, Message: "repository Error: filter value not set"}
	ErrResquestNotSet      = &echo.HTTPError{Internal: errors.New("repository Error: Request value not set"), Code: http.StatusBadRequest, Message: "repository Error: Request value not set"}
	ErrFailCreatingClient  = &echo.HTTPError{Internal: errors.New("repository Error: Failed to create DB client"), Code: http.StatusInternalServerError, Message: "repository Error: Failed to create DB client"}
	ErrIDnotFound          = &echo.HTTPError{Internal: errors.New("repository Error: Id not found"), Code: http.StatusBadRequest, Message: "repository Error: Id not found"}
	ErrDataTypeWrong       = &echo.HTTPError{Internal: errors.New("repository Error: Invalid argument passed"), Code: http.StatusBadRequest, Message: "repository Error: Invalid argument passed"}
	ErrInvalidStatus       = &echo.HTTPError{Internal: errors.New("invalid status value"), Code: http.StatusBadRequest, Message: "invalid status value"}
	ErrNotAuthenticated    = &echo.HTTPError{Internal: errors.New("authorization header is missing, please authenticate first"), Code: http.StatusUnauthorized, Message: "authorization header is missing, please authenticate first"}
	ErrUserIDNotMatch      = &echo.HTTPError{Internal: errors.New("user id not match with id in token"), Code: http.StatusUnauthorized, Message: "user id not match with id in token"}
	ErrAccountNotActive    = &echo.HTTPError{Internal: errors.New("account is not active"), Code: http.StatusBadRequest, Message: "account is not active"}
	ErrClientIDNotValid    = &echo.HTTPError{Internal: errors.New("client ID not valid"), Code: http.StatusBadRequest, Message: "client ID not valid"}
	ErrUserIDNotValid      = &echo.HTTPError{Internal: errors.New("user ID not valid"), Code: http.StatusBadRequest, Message: "user ID not valid"}
	ErrAgencyIDNotValid    = &echo.HTTPError{Internal: errors.New("agency ID not valid"), Code: http.StatusBadRequest, Message: "agency ID not valid"}
	ErrMissingCredentials  = &echo.HTTPError{Internal: errors.New("missing credentials"), Code: http.StatusBadRequest, Message: "missing credentials"}
	ValidStatus            = []Status{"active", "blocked"}
)

type StandartResponse struct {
	Message string `json:"message" xml:"message"`
}

type Claims struct {
	Id   string `json:"id" xml:"id"`
	Role string `json:"role" xml:"role"`
	jwt.RegisteredClaims
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
