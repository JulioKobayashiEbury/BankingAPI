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
}

func NewWithdrawalFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return withdrawalFirestore{
		databaseClient: dbClient,
	}
}

func (db withdrawalFirestore) Create(request interface{}) (interface{}, *model.Erro) {
	Withdrawal, ok := request.(*Withdrawal)
	if !ok {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":      Withdrawal.Account_id,
		"client_id":       Withdrawal.Client_id,
		"agency_id":       Withdrawal.Agency_id,
		"withdrawal":      Withdrawal.Withdrawal,
		"status":          true,
		"withdrawal_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	// Add withdrawal to account list
	return db.Get(&docRef.ID)
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
	Withdrawal := Withdrawal{}
	if err := docSnapshot.DataTo(&Withdrawal); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &Withdrawal, nil
}

func (db withdrawalFirestore) Update(request interface{}) *model.Erro {
	Withdrawal, ok := request.(*Withdrawal)
	if !ok {
		return model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"account_id":      Withdrawal.Account_id,
		"client_id":       Withdrawal.Client_id,
		"agency_id":       Withdrawal.Agency_id,
		"withdrawal":      Withdrawal.Withdrawal,
		"status":          true,
		"withdrawal_date": time.Now().Format(model.TimeLayout),
	}
	docRef := db.databaseClient.Collection(collection).Doc(Withdrawal.Withdrawal_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	log.Info().Msg("Account: " + Withdrawal.Withdrawal_id + " has been updated")

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
	WithdrawalSlice := make([]*Withdrawal, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		Withdrawal := &Withdrawal{}
		if err := docSnap.DataTo(&Withdrawal); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		Withdrawal.Withdrawal_id = docSnap.Ref.ID
		WithdrawalSlice = append(WithdrawalSlice, Withdrawal)
	}

	return &WithdrawalSlice, nil
}
