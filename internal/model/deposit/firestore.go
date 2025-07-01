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

func NewDepositFirestore(dbClient *firestore.Client) DepositRepository {
	return depositFirestore{
		databaseClient: dbClient,
	}
}

func (db depositFirestore) Create(ctx context.Context, request *Deposit) (*Deposit, *model.Erro) {
	entity := map[string]interface{}{
		"account_id":   request.Account_id,
		"client_id":    request.Client_id,
		"agency_id":    request.Agency_id,
		"user_id":      request.User_id,
		"deposit":      request.Deposit,
		"deposit_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return db.Get(ctx, &docRef.ID)
}

func (db depositFirestore) Delete(ctx context.Context, id *string) *model.Erro {
	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db depositFirestore) Get(ctx context.Context, id *string) (*Deposit, *model.Erro) {
	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + *id), HttpCode: http.StatusInternalServerError}
	}
	deposit := Deposit{}
	if err := docSnapshot.DataTo(&deposit); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	deposit.Deposit_id = docSnapshot.Ref.ID
	log.Info().Msg("Deposit: " + deposit.Deposit_id + " has been retrieved")
	return &deposit, nil
}

func (db depositFirestore) Update(ctx context.Context, request *Deposit) *model.Erro {
	entity := map[string]interface{}{
		"account_id":   request.Account_id,
		"client_id":    request.Client_id,
		"agency_id":    request.Agency_id,
		"user_id":      request.User_id,
		"deposit":      request.Deposit,
		"deposit_date": request.Deposit_date,
	}
	docRef := db.databaseClient.Collection(collection).Doc(request.Deposit_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	log.Info().Msg("Account: " + request.Deposit_id + " has been updated")

	return nil
}

func (db depositFirestore) GetAll(ctx context.Context) (*[]Deposit, *model.Erro) {
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

func (db depositFirestore) GetFilteredByID(ctx context.Context, filters *string) (*[]Deposit, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	query := db.databaseClient.Collection(collection).Query
	/*
		for _, filter := range *filters {
			token := model.TokenizeFilters(&filter)
			if len(*token) != 3 {
				return nil, model.InvalidFilterFormat
			}

			query = query.Where((*token)[0], (*token)[1], (*token)[2])
		}
	*/
	allDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	depositSlice := make([]Deposit, 0, len(allDocs))
	for _, docSnap := range allDocs {
		depositResponse := Deposit{}
		if err := docSnap.DataTo(&depositResponse); err != nil {
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}

		depositResponse.Deposit_id = docSnap.Ref.ID

		depositSlice = append(depositSlice, depositResponse)
	}

	return &depositSlice, nil
}
