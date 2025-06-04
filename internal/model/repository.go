package model

import (
	"context"
	"errors"
	"net/http"
	"os"

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
	ResquestNotSet     = &Erro{Err: errors.New("Request value not set"), HttpCode: http.StatusInternalServerError}
	FailCreatingClient = &Erro{Err: errors.New("Failed to create DB client"), HttpCode: http.StatusInternalServerError}
	IDnotFound         = &Erro{Err: errors.New("Id not founc"), HttpCode: http.StatusBadRequest}
	DataTypeWrong      = &Erro{Err: errors.New("Invalid argument passed"), HttpCode: http.StatusBadRequest}
)

type RepositoryInterface interface {
	Create(interface{}) (*string, *Erro)
	Delete(*string) *Erro
	Get(id *string) (interface{}, *Erro)
	Update(interface{}) *Erro
	GetAll() (interface{}, *Erro)
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
