package transfer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo"
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

func (db transferFirestore) Create(ctx context.Context, request *Transfer) (*Transfer, *echo.HTTPError) {
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
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return db.Get(ctx, &docRef.ID)
}

func (db transferFirestore) Delete(ctx context.Context, id *string) *echo.HTTPError {
	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return nil
}

func (db transferFirestore) Get(ctx context.Context, id *string) (*Transfer, *echo.HTTPError) {
	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &echo.HTTPError{Internal: errors.New("ID in collection: " + collection + " not found"), Code: http.StatusBadRequest, Message: fmt.Sprint("ID in collection: " + collection + " not found")}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &echo.HTTPError{Internal: errors.New("Nil account from snapshot" + (*id)), Code: http.StatusInternalServerError, Message: fmt.Sprint("Nil account from snapshot" + (*id))}
	}
	transferResponse := Transfer{}
	if err := docSnapshot.DataTo(&transferResponse); err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	transferResponse.Transfer_id = docSnapshot.Ref.ID

	return &transferResponse, nil
}

func (db transferFirestore) Update(ctx context.Context, request *Transfer) *echo.HTTPError {
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
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (db transferFirestore) GetAll(ctx context.Context) (*[]Transfer, *echo.HTTPError) {
	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	transferResponseSlice := make([]Transfer, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		transferResponse := Transfer{}
		if err := docSnap.DataTo(&transferResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		transferResponse.Transfer_id = docSnap.Ref.ID
		transferResponseSlice = append(transferResponseSlice, transferResponse)
	}
	return &transferResponseSlice, nil
}

func (db transferFirestore) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]Transfer, *echo.HTTPError) {
	if accountID == nil || len(*accountID) == 0 {
		return nil, model.ErrFilterNotSet
	}

	query := db.databaseClient.Collection(collection).Query.Where("account_id", "==", *accountID)

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
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	transferSlice := make([]Transfer, 0, len(allDocs))
	for _, docSnap := range allDocs {
		transferResponse := Transfer{}
		if err := docSnap.DataTo(&transferResponse); err != nil {
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}

		transferResponse.Transfer_id = docSnap.Ref.ID

		transferSlice = append(transferSlice, transferResponse)
	}
	return &transferSlice, nil
}
