package model

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

const (
	AccountsPath    = "accounts"
	UsersPath       = "users"
	ClientPath      = "clients"
	TransfersPath   = "transfers"
	AutoDebit       = "autodebit"
	AutoDebitLog    = "autodebitlog"
	DepositPath     = "deposits"
	WithdrawalsPath = "withdrawals"
)

var (
	InvalidFilterFormat = &Erro{Err: errors.New("repository Error: Invalid fitler format"), HttpCode: http.StatusBadRequest}
	FilterNotSet        = &Erro{Err: errors.New("repository Error: filter value not set"), HttpCode: http.StatusBadRequest}
	ResquestNotSet      = &Erro{Err: errors.New("repository Error: Request value not set"), HttpCode: http.StatusBadRequest}
	FailCreatingClient  = &Erro{Err: errors.New("repository Error: Failed to create DB client"), HttpCode: http.StatusInternalServerError}
	IDnotFound          = &Erro{Err: errors.New("repository Error: Id not founc"), HttpCode: http.StatusBadRequest}
	DataTypeWrong       = &Erro{Err: errors.New("repository Error: Invalid argument passed"), HttpCode: http.StatusBadRequest}
	InvalidStatus       = &Erro{Err: errors.New("invalid status value"), HttpCode: http.StatusBadRequest}

	ValidStatus = []Status{"active", "blocked"}
)

type RepositoryInterface interface {
	Create(interface{}) (interface{}, *Erro)
	Delete(*string) *Erro
	Get(id *string) (interface{}, *Erro)
	Update(interface{}) *Erro
	GetAll() (interface{}, *Erro)
	GetFiltered(*[]string) (interface{}, *Erro)
}

func GetFireStoreClient() (*firestore.Client, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		log.Error().Msg("GOOGLE_CLOUD_PROJECT environment variable not set.")
		return nil, errors.New("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Error().Msg("Failed to create client: %v")
		return nil, err
	}
	return client, nil
}

func TokenizeFilters(filters *string) *[]string {
	tokens := strings.Split(*filters, ",")
	return &tokens
}
