package repository

import (
	"errors"
	"fmt"
	"strconv"

	"BankingAPI/internal/initializers"

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

func GetFireStoreClient() (*firestore.Client, error) {
	return initializers.FireBaseApp.Firestore(initializers.Ctx)
}

func CreateObject(entity interface{}, collection string, createdID *uint32) error {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	defer clientDB.Close()
	docRef, _, err := clientDB.Collection(collection).Add(initializers.Ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	newID, err := strconv.ParseUint(docRef.ID, 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	(*createdID) = uint32(newID)
	return nil
}

func DeleteObject(objectID uint32, collection string) error {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return err
	}
	defer clientDB.Close()
	docRef := clientDB.Collection(collection).Doc(fmt.Sprintf("%v", objectID))
	_, err = docRef.Delete(initializers.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetTypeFromDB(typesID *uint32, collection string) (*firestore.DocumentSnapshot, error) {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		return nil, err
	}
	defer clientDB.Close()

	docRef := clientDB.Collection(collection).Doc(fmt.Sprintf("%v", *typesID))

	docSnapshot, err := docRef.Get(initializers.Ctx)
	if status.Code(err) == codes.NotFound {
		return nil, errors.New("Account, User or Client with this ID do not exists")
	}
	return docSnapshot, nil
}

func UpdateTypesDB(document *[]firestore.Update, typesID *uint32, collection string) error {
	clientDB, err := GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	defer clientDB.Close()

	docRef := clientDB.Collection((collection)).Doc(fmt.Sprintf("%v", *typesID))

	_, err = docRef.Update(initializers.Ctx, *document)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}
