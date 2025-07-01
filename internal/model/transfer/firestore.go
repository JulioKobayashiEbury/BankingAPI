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

func NewTransferFirestore(dbClient *firestore.Client) TransferRepository {
	return transferFirestore{
		databaseClient: dbClient,
	}
}

func (db transferFirestore) Create(ctx context.Context, request *Transfer) (*Transfer, *model.Erro) {
	entity := map[string]interface{}{
		"account_id":    request.Account_id,
		"account_to":    request.Account_to,
		"value":         request.Value,
		"user_id":       request.User_id,
		"user_to":       request.User_to,
		"register_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return db.Get(ctx, &docRef.ID)
}

func (db transferFirestore) Delete(ctx context.Context, id *string) *model.Erro {
	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db transferFirestore) Get(ctx context.Context, id *string) (*Transfer, *model.Erro) {
	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + (*id)), HttpCode: http.StatusInternalServerError}
	}
	transferResponse := Transfer{}
	if err := docSnapshot.DataTo(&transferResponse); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	transferResponse.Transfer_id = docSnapshot.Ref.ID

	return &transferResponse, nil
}

func (db transferFirestore) Update(ctx context.Context, request *Transfer) *model.Erro {
	entity := map[string]interface{}{
		"account_id":    request.Account_id,
		"account_to":    request.Account_to,
		"value":         request.Value,
		"register_date": request.Register_date,
	}
	docRef := db.databaseClient.Collection(collection).Doc(request.Transfer_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	return nil
}

func (db transferFirestore) GetAll(ctx context.Context) (*[]Transfer, *model.Erro) {
	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	transferResponseSlice := make([]Transfer, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		transferResponse := Transfer{}
		if err := docSnap.DataTo(&transferResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		transferResponse.Transfer_id = docSnap.Ref.ID
		transferResponseSlice = append(transferResponseSlice, transferResponse)
	}
	return &transferResponseSlice, nil
}

func (db transferFirestore) GetFilteredByID(ctx context.Context, filters *string) (*[]Transfer, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}

	query := db.databaseClient.Collection(collection).Query
	/* for _, filter := range *filters {
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

	transferSlice := make([]Transfer, 0, len(allDocs))
	for _, docSnap := range allDocs {
		transferResponse := Transfer{}
		if err := docSnap.DataTo(&transferResponse); err != nil {
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}

		transferResponse.Transfer_id = docSnap.Ref.ID

		transferSlice = append(transferSlice, transferResponse)
	}
	return &transferSlice, nil
}
