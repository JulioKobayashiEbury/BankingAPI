package account

import (
	"context"
	"errors"
	"net/http"
	"time"

	model "BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collection = "accounts"

var accountDatabase model.RepositoryInterface

type accountFirestore struct {
	databaseClient *firestore.Client
	updateList     map[string]interface{}
}

func NewAccountFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return accountFirestore{
		databaseClient: dbClient,
	}
}

func (db accountFirestore) AddUpdate(key string, value interface{}) {
	if db.updateList == nil {
		db.updateList = make(map[string]interface{})
	}
	db.updateList[key] = value
}

func (db accountFirestore) Create(request interface{}) (*string, *model.Erro) {
	_, accountRequest := interfaceToAccount(request)
	if accountRequest == nil {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"client_id":     accountRequest.Client_id,
		"user_id":       accountRequest.User_id,
		"agency_id":     accountRequest.Agency_id,
		"balance":       0.0,
		"register_date": time.Now().Format(model.TimeLayout),
		"status":        true,
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &docRef.ID, nil
}

func (db accountFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db accountFirestore) Get(id *string) (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + *id), HttpCode: http.StatusInternalServerError}
	}
	accountResponse := AccountResponse{}
	if err := docSnapshot.DataTo(&accountResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponse.Account_id = docSnapshot.Ref.ID
	return accountResponse, nil
}

func (db accountFirestore) Update(id *string) *model.Erro {
	if db.updateList == nil {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	updates := make([]firestore.Update, 0, 0)
	for key, value := range db.updateList {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection((collection)).Doc(*id)

	docSnap, _ := docRef.Get(ctx)
	if !docSnap.Exists() {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID from collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db accountFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponseSlice := make([]*AccountResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		accountResponse := &AccountResponse{}
		if err := docSnap.DataTo(&accountResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		accountResponse.Account_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		accountResponseSlice = append(accountResponseSlice, accountResponse)
	}
	return &accountResponseSlice, nil
}

func interfaceToAccount(argument interface{}) (*AccountResponse, *AccountRequest) {
	if obj, ok := argument.(AccountResponse); ok {
		return &obj, nil
	}
	if obj, ok := argument.(AccountRequest); ok {
		return nil, &obj
	}
	return nil, nil
}
