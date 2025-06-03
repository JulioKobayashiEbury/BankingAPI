package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collection = "users"

type userFirestore struct {
	databaseClient *firestore.Client
	updateList     map[string]interface{}
}

func NewUserFireStore(dbClient *firestore.Client) model.RepositoryInterface {
	return userFirestore{
		databaseClient: dbClient,
	}
}

func (db userFirestore) AddUpdate(key string, value interface{}) {
	if db.updateList == nil {
		db.updateList = make(map[string]interface{})
	}
}

func (db userFirestore) Create(request interface{}) (*string, *model.Erro) {
	userRequest, _ := interfaceToUser(request)
	if userRequest == nil {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return nil, model.FailCreatingClient
	}
	defer clientDB.Close()
	entity := map[string]interface{}{
		"name":          userRequest.Name,
		"document":      userRequest.Document,
		"password":      userRequest.Password,
		"register_date": time.Now().Format(model.TimeLayout),
		"status":        true,
	}
	docRef, _, err := clientDB.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &docRef.ID, nil
}

func (db userFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()

	docRef := clientDB.Collection(collection).Doc(*id)
	_, err = docRef.Delete(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db userFirestore) Get(id *string) (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docSnapshot, err := clientDB.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, model.IDnotFound
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + (*id)), HttpCode: http.StatusInternalServerError}
	}
	userResponse := User{}
	if err := docSnapshot.DataTo(&userResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	userResponse.User_id = docSnapshot.Ref.ID
	return &userResponse, nil
}

func (db userFirestore) Update(id *string) *model.Erro {
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

	clientDB, err := model.GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docRef := clientDB.Collection((collection)).Doc(*id)

	docSnap, _ := docRef.Get(ctx)
	if !docSnap.Exists() {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID from collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	_, err = docRef.Update(ctx, updates)
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db userFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	clientDB, err := model.GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	iterator := clientDB.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	userResponseSlice := make([]*UserResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		userResponse := &UserResponse{}
		if err := docSnap.DataTo(&userResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		userResponse.User_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		userResponseSlice = append(userResponseSlice, userResponse)
	}
	return &userResponseSlice, nil
}

func interfaceToUser(argument interface{}) (*UserRequest, *UserResponse) {
	if obj, ok := argument.(UserRequest); ok {
		return &obj, nil
	}
	if obj, ok := argument.(UserResponse); ok {
		return nil, &obj
	}
	return nil, nil
}
