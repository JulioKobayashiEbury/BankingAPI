package model

import (
	"context"
	"errors"
	"net/http"
	"os"

	model "BankingAPI/internal/model/types"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

var Ctx context.Context

func GetFireStoreClient() (*firestore.Client, error) {
	Ctx = context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		log.Error().Msg("GOOGLE_CLOUD_PROJECT environment variable not set.")
		return nil, errors.New("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	client, err := firestore.NewClient(Ctx, projectID)
	if err != nil {
		log.Error().Msg("Failed to create client: %v")
		return nil, err
	}
	Ctx.Done()
	return client, nil
}

func CreateObject(entity *map[string]interface{}, collection string, createdID *string) *model.Erro {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()
	docRef, _, err := clientDB.Collection(collection).Add(Ctx, (*entity))
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	Ctx.Done()
	(*createdID) = docRef.ID
	return nil
}

func DeleteObject(objectID *string, collection string) *model.Erro {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()
	docRef := clientDB.Collection(collection).Doc(*objectID)
	_, err = docRef.Delete(Ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	Ctx.Done()
	return nil
}

func GetTypeFromDB(typesID *string, collection string) (*firestore.DocumentSnapshot, *model.Erro) {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docSnapshot, err := clientDB.Collection(collection).Doc((*typesID)).Get(Ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + (*typesID))
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + (*typesID)), HttpCode: http.StatusInternalServerError}
	}
	Ctx.Done()
	return docSnapshot, nil
}

func UpdateTypesDB(document *[]firestore.Update, typesID *string, collection string) *model.Erro {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docRef := clientDB.Collection((collection)).Doc(*typesID)

	docSnap, _ := docRef.Get(Ctx)
	if !docSnap.Exists() {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID from collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	_, err = docRef.Update(Ctx, *document)
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	Ctx.Done()
	return nil
}

func GetAllByTypeDB(collection string) ([]*firestore.DocumentSnapshot, *model.Erro) {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	iterator := clientDB.Collection(collection).Documents(Ctx)

	all, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	Ctx.Done()
	return all, nil
}
