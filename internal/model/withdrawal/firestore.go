package withdrawal

import (
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

type WithdrawalFirestore struct {
	Request  *WithdrawalRequest
	Response *WithdrawalResponse
	Slice    *[]*WithdrawalResponse
	model.Repository
}

func (db *WithdrawalFirestore) Create() *model.Erro {
	if db.Request == nil {
		log.Error().Msg("Client Request not set for DB creation")
		return model.ResquestNotSet
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()
	entity := map[string]interface{}{
		"account_id":      db.Request.Account_id,
		"client_id":       db.Request.Client_id,
		"agency_id":       db.Request.Agency_id,
		"withdrawal":      db.Request.Withdrawal,
		"status":          true,
		"withdrawal_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := clientDB.Collection(collection).Add(*ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	db.Response.Withdrawal_id = docRef.ID
	// Add withdrawal to account list
	return nil
}

func (db *WithdrawalFirestore) Delete() *model.Erro {
	if db.Request.Client_id == "" {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()
	docRef := clientDB.Collection(collection).Doc(db.Request.Withdrawal_id)
	_, err = docRef.Delete(*ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *WithdrawalFirestore) Get() *model.Erro {
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docSnapshot, err := clientDB.Collection(collection).Doc(db.Request.Withdrawal_id).Get(*ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + db.Request.Withdrawal_id)
		return &model.Erro{Err: errors.New("Nil account from snapshot" + (db.Request.Withdrawal_id)), HttpCode: http.StatusInternalServerError}
	}
	if err := docSnapshot.DataTo(db.Response); err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *WithdrawalFirestore) Update() *model.Erro {
	if db.GetUpdateList() == nil {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	updates := make([]firestore.Update, 0, 0)
	for key, value := range *db.GetUpdateList() {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docRef := clientDB.Collection((collection)).Doc(db.Request.Withdrawal_id)

	docSnap, _ := docRef.Get(*ctx)
	if !docSnap.Exists() {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID from collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	_, err = docRef.Update(*ctx, updates)
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *WithdrawalFirestore) GetAll() *model.Erro {
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	iterator := clientDB.Collection(collection).Documents(*ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	withdrawalResponseSlice := make([]*WithdrawalResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		withdrawalResponse := &WithdrawalResponse{}
		if err := docSnap.DataTo(&withdrawalResponse); err != nil {
			log.Error().Msg(err.Error())
			return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		withdrawalResponse.Withdrawal_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		withdrawalResponseSlice = append(withdrawalResponseSlice, withdrawalResponse)
	}
	db.Slice = &withdrawalResponseSlice
	return nil
}
