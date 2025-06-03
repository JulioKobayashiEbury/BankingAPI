package withdrawal

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

const collection = "withdrawal"

type withdrawalFirestore struct {
	databaseClient *firestore.Client
	updateList     map[string]interface{}
}

func (db withdrawalFirestore) AddUpdate(key string, value interface{}) {
	if db.updateList == nil {
		db.updateList = make(map[string]interface{})
	}
}

func (db withdrawalFirestore) Create(request interface{}) (*string, *model.Erro) {
	withdrawalRequest, ok := request.(WithdrawalRequest)
	if !ok {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":      withdrawalRequest.Account_id,
		"client_id":       withdrawalRequest.Client_id,
		"agency_id":       withdrawalRequest.Agency_id,
		"withdrawal":      withdrawalRequest.Withdrawal,
		"status":          true,
		"withdrawal_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	// Add withdrawal to account list
	return &docRef.ID, nil
}

func (db withdrawalFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()
	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db withdrawalFirestore) Get(id *string) (interface{}, *model.Erro) {
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
	withdrawalResponse := WithdrawalResponse{}
	if err := docSnapshot.DataTo(&withdrawalResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &withdrawalResponse, nil
}

func (db withdrawalFirestore) Update(id *string) *model.Erro {
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

func (db withdrawalFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	withdrawalResponseSlice := make([]*WithdrawalResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		withdrawalResponse := &WithdrawalResponse{}
		if err := docSnap.DataTo(&withdrawalResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		withdrawalResponse.Withdrawal_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		withdrawalResponseSlice = append(withdrawalResponseSlice, withdrawalResponse)
	}

	return &withdrawalResponseSlice, nil
}
