package account

import (
	"context"
	"errors"
	"net/http"
	"time"

	model "BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collection = "accounts"

type accountFirestore struct {
	databaseClient *firestore.Client
}

func NewAccountFirestore(dbClient *firestore.Client) model.RepositoryInterface {
	return accountFirestore{
		databaseClient: dbClient,
	}
}

func (db accountFirestore) Create(request interface{}) (interface{}, *model.Erro) {
	accountRequest, ok := request.(*Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"client_id":     accountRequest.Client_id,
		"user_id":       accountRequest.User_id,
		"agency_id":     accountRequest.Agency_id,
		"balance":       0.0,
		"register_date": time.Now().Format(model.TimeLayout),
		"status":        model.ValidStatus[0],
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return db.Get(&docRef.ID)
}

func (db accountFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db accountFirestore) Get(id *string) (interface{}, *model.Erro) {
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
	accountResponse := Account{}
	if err := docSnapshot.DataTo(&accountResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponse.Account_id = docSnapshot.Ref.ID
	return &accountResponse, nil
}

func (db accountFirestore) Update(request interface{}) *model.Erro {
	accountRequest, ok := request.(*Account)
	if !ok {
		return model.DataTypeWrong
	}

	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"client_id":     accountRequest.Client_id,
		"user_id":       accountRequest.User_id,
		"agency_id":     accountRequest.Agency_id,
		"balance":       accountRequest.Balance,
		"register_date": accountRequest.Register_date,
		"status":        accountRequest.Status,
	}
	docRef := db.databaseClient.Collection(collection).Doc(accountRequest.Account_id)

	if _, err := docRef.Set(ctx, entity); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	log.Info().Msg("Account: " + accountRequest.Account_id + " has been updated")

	return nil
}

func (db accountFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponseSlice := make([]Account, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		accountResponse := Account{}
		if err := docSnap.DataTo(&accountResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		accountResponse.Account_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		accountResponseSlice = append(accountResponseSlice, accountResponse)
	}
	return &accountResponseSlice, nil
}

func (db accountFirestore) GetFiltered(filters *[]string) (interface{}, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	ctx := context.Background()
	defer ctx.Done()

	query := db.databaseClient.Collection(collection).Query

	for _, filter := range *filters {
		tokens := model.TokenizeFilters(&filter)
		if len(*tokens) != 3 {
			log.Error().Msg("Invalid filter format: " + filter)
			return nil, model.InvalidFilterFormat
		}
		query = query.Where((*tokens)[0], (*tokens)[1], (*tokens)[2])
	}

	allDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	accountResponseSlice := make([]Account, 0, len(allDocs))
	for _, docSnap := range allDocs {
		accountResponse := Account{}
		if err := docSnap.DataTo(&accountResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		accountResponse.Account_id = docSnap.Ref.ID

		accountResponseSlice = append(accountResponseSlice, accountResponse)

	}

	return &accountResponseSlice, nil
}
