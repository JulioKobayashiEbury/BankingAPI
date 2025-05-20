package model

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AccountsPath = "accounts"
	UsersPath    = "users"
	ClientPath   = "clients"
)

var Ctx context.Context

func GetFireStoreClient() (*firestore.Client, error) {
	Ctx = context.Background()
	os.Setenv("FIRESTORE_EMULATOR_HOST", "0.0.0.0:8080")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "banking")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Error().Msg("GOOGLE_CLOUD_PROJECT environment variable not set.")
		return nil, errors.New("GOOGLE_CLOUD_PROJECT environment variable not set.")
	}

	client, err := firestore.NewClient(Ctx, projectID)
	if err != nil {
		log.Error().Msg("Failed to create client: %v")
		return nil, err
	}
	return client, nil
}

func CreateObject(entity *map[string]interface{}, collection string, createdID *string) error {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	defer clientDB.Close()
	docRef, _, err := clientDB.Collection(collection).Add(Ctx, (*entity))
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	(*createdID) = docRef.ID
	return nil
}

func DeleteObject(objectID *string, collection string) error {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return err
	}
	defer clientDB.Close()
	docRef := clientDB.Collection(collection).Doc(fmt.Sprintf("%v", objectID))
	_, err = docRef.Delete(Ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetTypeFromDB(typesID *string, collection string) (*firestore.DocumentSnapshot, error) {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return nil, err
	}
	defer clientDB.Close()

	docRef := clientDB.Collection(collection).Doc(fmt.Sprintf("%v", *typesID))

	docSnapshot, err := docRef.Get(Ctx)
	if status.Code(err) == codes.NotFound {
		return nil, errors.New("Account, User or Client with this ID do not exists")
	}
	if docSnapshot == nil {
		return nil, errors.New("Nil account")
	}
	return docSnapshot, nil
}

func UpdateTypesDB(document *[]firestore.Update, typesID *string, collection string) error {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	defer clientDB.Close()

	docRef := clientDB.Collection((collection)).Doc(fmt.Sprintf("%v", *typesID))

	_, err = docRef.Update(Ctx, *document)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}
