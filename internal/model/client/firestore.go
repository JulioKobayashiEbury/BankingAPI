package client

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

const collection = "clients"

type clientFirestore struct {
	databaseClient *firestore.Client
	updateList     map[string]interface{}
}

func NewClientFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return clientFirestore{
		databaseClient: dbClient,
	}
}

func (db clientFirestore) AddUpdate(key string, value interface{}) {
	if db.updateList == nil {
		db.updateList = make(map[string]interface{})
	}
	db.updateList[key] = value
}

func (db clientFirestore) Create(request interface{}) (*string, *model.Erro) {
	clientRequest, _ := interfaceToClient(request)
	if clientRequest == nil {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"user_id":       clientRequest.User_id,
		"name":          clientRequest.Name,
		"document":      clientRequest.Document,
		"status":        true,
		"register_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &docRef.ID, nil
}

func (db clientFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection(collection).Doc(*id)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db clientFirestore) Get(id *string) (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, model.IDnotFound
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + (*id)), HttpCode: http.StatusInternalServerError}
	}
	clientResponse := ClientResponse{}
	if err := docSnapshot.DataTo(&clientResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	clientResponse.Client_id = docSnapshot.Ref.ID
	return &clientResponse, nil
}

func (db clientFirestore) Update(id *string) *model.Erro {
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

func (db clientFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	clientResponseSlice := make([]*ClientResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		clientResponse := &ClientResponse{}
		if err := docSnap.DataTo(&clientResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		clientResponse.Client_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		clientResponseSlice = append(clientResponseSlice, clientResponse)
	}
	return &clientResponseSlice, nil
}

func interfaceToClient(argument interface{}) (*ClientRequest, *ClientResponse) {
	if obj, ok := argument.(ClientRequest); ok {
		return &obj, nil
	}
	if obj, ok := argument.(ClientResponse); ok {
		return nil, &obj
	}
	return nil, nil
}
