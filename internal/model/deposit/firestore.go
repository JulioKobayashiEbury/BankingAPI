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
}

func NewDepositFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return depositFirestore{
		databaseClient: dbClient,
	}
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
	Deposit := Deposit{}
	if err := docSnapshot.DataTo(&Deposit); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &Deposit, nil
}

func (db depositFirestore) Update(request interface{}) *model.Erro {
	depositRequest, ok := request.(*Deposit)
	if !ok {
		return model.DataTypeWrong
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
	docRef := db.databaseClient.Collection(collection).Doc(depositRequest.Deposit_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	log.Info().Msg("Account: " + depositRequest.Deposit_id + " has been updated")

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
	DepositSlice := make([]Deposit, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		depositReponse := Deposit{}
		if err := docSnap.DataTo(&depositReponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		depositReponse.Deposit_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		DepositSlice = append(DepositSlice, depositReponse)
	}
	return &DepositSlice, nil
}

func interfaceToDeposit(argument interface{}) (*Deposit, *Deposit) {
	if obj, ok := argument.(Deposit); ok {
		return &obj, nil
	}
	if obj, ok := argument.(Deposit); ok {
		return nil, &obj
	}
	return nil, nil
}
