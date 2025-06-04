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
}

func NewClientFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return clientFirestore{
		databaseClient: dbClient,
	}
}

func (db clientFirestore) Create(request interface{}) (*string, *model.Erro) {
	clientRequest, ok := request.(ClientRequest)
	if !ok {
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

func (db clientFirestore) Update(request interface{}) *model.Erro {
	clientRequest, ok := request.(*ClientRequest)
	if !ok {
		return model.DataTypeWrong
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
	docRef := db.databaseClient.Collection(collection).Doc(clientRequest.Client_id)

	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	log.Info().Msg("Account: " + clientRequest.Client_id + " has been updated")

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
