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
)

type RepositoryInterface interface {
	Create() *Erro
	Delete() *Erro
	Get() *Erro
	Update() *Erro
	GetAll() *Erro
}

type Repository struct {
	updateList map[string]interface{}
	RepositoryInterface
}

func GetFireStoreClient() (*context.Context, *firestore.Client, error) {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		log.Error().Msg("GOOGLE_CLOUD_PROJECT environment variable not set.")
		return nil, nil, errors.New("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Error().Msg("Failed to create client: %v")
		return nil, nil, err
	}
	ctx.Done()
	return &ctx, client, nil
}

func (db *Repository) AddUpdate(key string, value interface{}) {
	if (*db).updateList == nil {
		(*db).updateList = make(map[string]interface{})
	}
	(*db).updateList[key] = value
}

func (db *Repository) GetUpdateList() *map[string]interface{} {
	return &db.updateList
}
