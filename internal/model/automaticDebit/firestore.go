package automaticdebit

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

const collection = "autodebit"

type AutoDebitFirestore struct {
	AutoDebit *AutomaticDebit
	Slice     *[]*AutomaticDebit
	model.Repository
}

func (db *AutoDebitFirestore) Create() *model.Erro {
	if db.AutoDebit == nil {
		log.Error().Msg("Client Request not set for DB creation")
		return model.ResquestNotSet
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()
	entity := map[string]interface{}{
		"account_id":      db.AutoDebit.Account_id,
		"client_id":       db.AutoDebit.Client_id,
		"agency_id":       db.AutoDebit.Agency_id,
		"value":           db.AutoDebit.Value,
		"debit_day":       db.AutoDebit.Debit_day,
		"expiration_date": db.AutoDebit.Expiration_date,
		"status":          true,
		"register_date":   time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := clientDB.Collection(collection).Add(*ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	db.AutoDebit = &AutomaticDebit{
		Debit_id: docRef.ID,
	}
	return nil
}

func (db *AutoDebitFirestore) Delete() *model.Erro {
	if db.AutoDebit.Debit_id == "" {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()
	docRef := clientDB.Collection(collection).Doc(db.AutoDebit.Debit_id)
	_, err = docRef.Delete(*ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *AutoDebitFirestore) Get() *model.Erro {
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docSnapshot, err := clientDB.Collection(collection).Doc(db.AutoDebit.Debit_id).Get(*ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + db.AutoDebit.Debit_id)
		return &model.Erro{Err: errors.New("Nil account from snapshot" + (db.AutoDebit.Debit_id)), HttpCode: http.StatusInternalServerError}
	}
	if err := docSnapshot.DataTo(db.AutoDebit); err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *AutoDebitFirestore) Update() *model.Erro {
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

	docRef := clientDB.Collection((collection)).Doc(db.AutoDebit.Debit_id)

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

func (db *AutoDebitFirestore) GetAll() *model.Erro {
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
	autoebitResponseSlice := make([]*AutomaticDebit, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		autodebitReponse := &AutomaticDebit{}
		if err := docSnap.DataTo(&autodebitReponse); err != nil {
			log.Error().Msg(err.Error())
			return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		autodebitReponse.Debit_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		autoebitResponseSlice = append(autoebitResponseSlice, autodebitReponse)
	}
	db.Slice = &autoebitResponseSlice
	return nil
}
