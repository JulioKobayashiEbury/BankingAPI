package automaticdebit

import (
	"context"
	"errors"
	"net/http"
	"time"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collection = "autodebit"

type autoDebitFirestore struct {
	databaseClient *firestore.Client
}

func NewAutoDebitFirestore(dbClient *firestore.Client) AutoDebitRepository {
	return autoDebitFirestore{
		databaseClient: dbClient,
	}
}

func (db autoDebitFirestore) Create(ctx context.Context, request *AutomaticDebit) (*AutomaticDebit, *echo.HTTPError) {
	entity := map[string]interface{}{
		"account_id":      request.Account_id,
		"user_id":         request.User_id,
		"agency_id":       request.Agency_id,
		"value":           request.Value,
		"debit_day":       request.Debit_day,
		"expiration_date": request.Expiration_date,
		"status":          model.ValidStatus[0],
		"register_date":   time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError}
	}
	return db.Get(ctx, &docRef.ID)
}

func (db autoDebitFirestore) Delete(ctx context.Context, id *string) *echo.HTTPError {
	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError}
	}
	return nil
}

func (db autoDebitFirestore) Get(ctx context.Context, id *string) (*AutomaticDebit, *echo.HTTPError) {
	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &echo.HTTPError{Internal: errors.New("ID in collection: " + collection + " not found"), Code: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &echo.HTTPError{Internal: errors.New("Nil account from snapshot" + (*id)), Code: http.StatusInternalServerError}
	}
	autoDebitResponse := AutomaticDebit{}
	if err := docSnapshot.DataTo(&autoDebitResponse); err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError}
	}
	autoDebitResponse.Debit_id = docSnapshot.Ref.ID
	return &autoDebitResponse, nil
}

func (db autoDebitFirestore) Update(ctx context.Context, request *AutomaticDebit) *echo.HTTPError {
	entity := map[string]interface{}{
		"account_id":      request.Account_id,
		"user_id":         request.User_id,
		"agency_id":       request.Agency_id,
		"value":           request.Value,
		"debit_day":       request.Debit_day,
		"expiration_date": request.Expiration_date,
		"status":          true,
		"register_date":   request.Register_date,
	}
	docRef := db.databaseClient.Collection(collection).Doc(request.Debit_id)

	if _, err := docRef.Set(ctx, entity); err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	log.Info().Msg("Account: " + request.Debit_id + " has been updated")

	return nil
}

func (db autoDebitFirestore) GetAll(ctx context.Context) (*[]AutomaticDebit, *echo.HTTPError) {
	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	autoebitResponseSlice := make([]AutomaticDebit, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		autodebitReponse := AutomaticDebit{}
		if err := docSnap.DataTo(&autodebitReponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		autodebitReponse.Debit_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		autoebitResponseSlice = append(autoebitResponseSlice, autodebitReponse)
	}
	return &autoebitResponseSlice, nil
}

func (db autoDebitFirestore) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]AutomaticDebit, *echo.HTTPError) {
	if accountID == nil || len(*accountID) == 0 {
		return nil, model.ErrFilterNotSet
	}

	query := db.databaseClient.Collection(collection).Query

	query = query.Where("account_id", "==", *accountID)

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
	autodebitResponseSlice := make([]AutomaticDebit, 0, len(allDocs))
	for _, docSnap := range allDocs {
		autodebitResponse := AutomaticDebit{}
		if err := docSnap.DataTo(&autodebitResponse); err != nil {
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		autodebitResponse.Debit_id = docSnap.Ref.ID

		autodebitResponseSlice = append(autodebitResponseSlice, autodebitResponse)
	}
	return &autodebitResponseSlice, nil
}
