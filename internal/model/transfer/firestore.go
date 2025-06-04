package transfer

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

const collection = "transfers"

type transferFirestore struct {
	databaseClient *firestore.Client
}

func NewTransferFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return transferFirestore{
		databaseClient: dbClient,
	}
}

func (db transferFirestore) Create(request interface{}) (*string, *model.Erro) {
	transferRequest, ok := request.(TransferRequest)
	if !ok {
		return nil, model.DataTypeWrong
	}

	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":    transferRequest.Account_id,
		"account_to":    transferRequest.Account_to,
		"value":         transferRequest.Value,
		"register_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &docRef.ID, nil
}

func (db transferFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db transferFirestore) Get(id *string) (interface{}, *model.Erro) {
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
	transferResponse := TransferResponse{}
	if err := docSnapshot.DataTo(&transferResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &transferResponse, nil
}

func (db transferFirestore) Update(request interface{}) *model.Erro {
	transferRequest, ok := request.(*TransferRequest)
	if !ok {
		return model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":    transferRequest.Account_id,
		"account_to":    transferRequest.Account_to,
		"value":         transferRequest.Value,
		"register_date": time.Now().Format(model.TimeLayout),
	}
	docRef := db.databaseClient.Collection(collection).Doc(transferRequest.Transfer_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	return nil
}

func (db transferFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	transferResponseSlice := make([]*TransferResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		transferResponse := &TransferResponse{}
		if err := docSnap.DataTo(&transferResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		transferResponse.Transfer_id = docSnap.Ref.ID
		transferResponseSlice = append(transferResponseSlice, transferResponse)
	}
	return &transferResponseSlice, nil
}
