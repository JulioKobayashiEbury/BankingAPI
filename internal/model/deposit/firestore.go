package deposit

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

const collection = "deposit"

type depositFirestore struct {
	databaseClient *firestore.Client
	updateList     map[string]interface{}
}

func NewDepositFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return depositFirestore{
		databaseClient: dbClient,
	}
}

func (db depositFirestore) AddUpdate(key string, value interface{}) {
	if db.updateList == nil {
		db.updateList = make(map[string]interface{})
	}
	db.updateList[key] = value
}

func (db depositFirestore) Create(request interface{}) (*string, *model.Erro) {
	depositRequest, _ := interfaceToDeposit(request)
	if depositRequest == nil {
		return nil, model.DataTypeWrong
	}

	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":      depositRequest.Account_id,
		"client_id":       depositRequest.Client_id,
		"agency_id":       depositRequest.Agency_id,
		"deposit":         depositRequest.Deposit,
		"status":          true,
		"withdrawal_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &docRef.ID, nil
}

func (db depositFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection(collection).Doc(*id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db depositFirestore) Get(id *string) (interface{}, *model.Erro) {
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
	depositResponse := DepositResponse{}
	if err := docSnapshot.DataTo(&depositResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &depositResponse, nil
}

func (db depositFirestore) Update(id *string) *model.Erro {
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
	_, err := docRef.Update(ctx, updates)
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db depositFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	depositResponseSlice := make([]*DepositResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		depositReponse := &DepositResponse{}
		if err := docSnap.DataTo(&depositReponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		depositReponse.Deposit_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		depositResponseSlice = append(depositResponseSlice, depositReponse)
	}
	return &depositResponseSlice, nil
}

func interfaceToDeposit(argument interface{}) (*DepositRequest, *DepositResponse) {
	if obj, ok := argument.(DepositRequest); ok {
		return &obj, nil
	}
	if obj, ok := argument.(DepositResponse); ok {
		return nil, &obj
	}
	return nil, nil
}
