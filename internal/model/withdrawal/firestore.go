package withdrawal

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

const collection = "withdrawal"

type withdrawalFirestore struct {
	databaseClient *firestore.Client
}

func NewWithdrawalFirestore(dbClient *firestore.Client) WithdrawalRepository {
	return withdrawalFirestore{
		databaseClient: dbClient,
	}
}

func (db withdrawalFirestore) Create(ctx context.Context, request *Withdrawal) (*Withdrawal, *echo.HTTPError) {
	entity := map[string]interface{}{
		"account_id":      request.Account_id,
		"user_id":         request.User_id,
		"agency_id":       request.Agency_id,
		"withdrawal":      request.Withdrawal,
		"withdrawal_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	// Add withdrawal to account list
	return db.Get(ctx, &docRef.ID)
}

func (db withdrawalFirestore) Delete(ctx context.Context, id *string) *echo.HTTPError {
	docRef := db.databaseClient.Collection(collection).Doc(*id)

	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (db withdrawalFirestore) Get(ctx context.Context, id *string) (*Withdrawal, *echo.HTTPError) {
	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &echo.HTTPError{Internal: errors.New("ID in collection: " + collection + " not found"), Code: http.StatusBadRequest, Message: fmt.Sprint("ID in collection: " + collection + " not found")}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &echo.HTTPError{Internal: errors.New("Nil account from snapshot" + (*id)), Code: http.StatusInternalServerError, Message: fmt.Sprint("Nil account from snapshot" + (*id))}
	}
	withdrawal := Withdrawal{}
	if err := docSnapshot.DataTo(&withdrawal); err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	withdrawal.Withdrawal_id = docSnapshot.Ref.ID

	return &withdrawal, nil
}

func (db withdrawalFirestore) Update(ctx context.Context, request *Withdrawal) *echo.HTTPError {
	entity := map[string]interface{}{
		"account_id":      request.Account_id,
		"user_id":         request.User_id,
		"agency_id":       request.Agency_id,
		"withdrawal":      request.Withdrawal,
		"status":          true,
		"withdrawal_date": request.Withdrawal_date,
	}
	docRef := db.databaseClient.Collection(collection).Doc(request.Withdrawal_id)
	_, err := docRef.Set(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	log.Info().Msg("Account: " + request.Withdrawal_id + " has been updated")

	return nil
}

func (db withdrawalFirestore) GetAll(ctx context.Context) (*[]Withdrawal, *echo.HTTPError) {
	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	WithdrawalSlice := make([]Withdrawal, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		Withdrawal := Withdrawal{}
		if err := docSnap.DataTo(&Withdrawal); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		Withdrawal.Withdrawal_id = docSnap.Ref.ID
		WithdrawalSlice = append(WithdrawalSlice, Withdrawal)
	}

	return &WithdrawalSlice, nil
}

func (db withdrawalFirestore) GetFilteredByAccountID(ctx context.Context, accountID *string) (*[]Withdrawal, *echo.HTTPError) {
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

	withdrawalsSlice := make([]Withdrawal, 0, len(allDocs))
	for _, docSnap := range allDocs {
		withdrawal := Withdrawal{}
		if err := docSnap.DataTo(&withdrawal); err != nil {
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}

		withdrawal.Withdrawal_id = docSnap.Ref.ID

		withdrawalsSlice = append(withdrawalsSlice, withdrawal)
	}

	return &withdrawalsSlice, nil
}
