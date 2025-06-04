package automaticdebit

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

const collection = "autodebit"

type autoDebitFirestore struct {
	databaseClient *firestore.Client
}

func NewAutoDebitFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return autoDebitFirestore{
		databaseClient: dbClient,
	}
}

func (db autoDebitFirestore) Create(request interface{}) (*string, *model.Erro) {
	autoDebitRequest, ok := request.(AutomaticDebitRequest)
	if !ok {
		return nil, model.DataTypeWrong
	}

	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":      autoDebitRequest.Account_id,
		"client_id":       autoDebitRequest.Client_id,
		"agency_id":       autoDebitRequest.Agency_id,
		"value":           autoDebitRequest.Value,
		"debit_day":       autoDebitRequest.Debit_day,
		"expiration_date": autoDebitRequest.Expiration_date,
		"status":          true,
		"register_date":   time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &docRef.ID, nil
}

func (db autoDebitFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db autoDebitFirestore) Get(id *string) (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + (*id)), HttpCode: http.StatusInternalServerError}
	}
	autoDebitResponse := AutomaticDebitResponse{}
	if err := docSnapshot.DataTo(&autoDebitResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	autoDebitResponse.Debit_id = docSnapshot.Ref.ID
	return &autoDebitResponse, nil
}

func (db autoDebitFirestore) Update(request interface{}) *model.Erro {
	autoDebitRequest, ok := request.(*AutomaticDebitRequest)
	if !ok {
		return model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":      autoDebitRequest.Account_id,
		"client_id":       autoDebitRequest.Client_id,
		"agency_id":       autoDebitRequest.Agency_id,
		"value":           autoDebitRequest.Value,
		"debit_day":       autoDebitRequest.Debit_day,
		"expiration_date": autoDebitRequest.Expiration_date,
		"status":          true,
		"register_date":   time.Now().Format(model.TimeLayout),
	}
	docRef := db.databaseClient.Collection(collection).Doc(autoDebitRequest.Debit_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	log.Info().Msg("Account: " + autoDebitRequest.Debit_id + " has been updated")

	return nil
}

func (db autoDebitFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	autoebitResponseSlice := make([]*AutomaticDebitResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		autodebitReponse := &AutomaticDebitResponse{}
		if err := docSnap.DataTo(&autodebitReponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		autodebitReponse.Debit_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		autoebitResponseSlice = append(autoebitResponseSlice, autodebitReponse)
	}
	return &autoebitResponseSlice, nil
}
