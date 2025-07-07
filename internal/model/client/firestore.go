package client

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

const collection = "clients"

type clientFirestore struct {
	databaseclient *firestore.Client
}

func NewClientFirestore(dbclient *firestore.Client) ClientRepository {
	return clientFirestore{
		databaseclient: dbclient,
	}
}

func (db clientFirestore) Create(ctx context.Context, request *Client) (*Client, *echo.HTTPError) {
	entity := map[string]interface{}{
		"user_id":       request.User_id,
		"name":          request.Name,
		"document":      request.Document,
		"register_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseclient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return db.Get(ctx, &docRef.ID)
}

func (db clientFirestore) Delete(ctx context.Context, id *string) *echo.HTTPError {
	docRef := db.databaseclient.Collection(collection).Doc(*id)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return nil
}

func (db clientFirestore) Get(ctx context.Context, id *string) (*Client, *echo.HTTPError) {
	docSnapshot, err := db.databaseclient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, model.ErrIDnotFound
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &echo.HTTPError{Internal: errors.New("Nil account from snapshot" + (*id)), Code: http.StatusInternalServerError, Message: fmt.Sprint("Nil account from snapshot" + (*id))}
	}
	client := Client{}
	if err := docSnapshot.DataTo(&client); err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	client.Client_id = docSnapshot.Ref.ID
	return &client, nil
}

func (db clientFirestore) Update(ctx context.Context, request *Client) *echo.HTTPError {
	entity := map[string]interface{}{
		"user_id":       request.User_id,
		"name":          request.Name,
		"document":      request.Document,
		"register_date": request.Register_date,
	}
	docRef := db.databaseclient.Collection(collection).Doc(request.Client_id)

	if _, err := docRef.Set(ctx, entity); err != nil {
		log.Error().Msg(err.Error())
		return &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	log.Info().Msg("Account: " + request.Client_id + " has been updated")

	return nil
}

func (db clientFirestore) GetAll(ctx context.Context) (*[]Client, *echo.HTTPError) {
	iterator := db.databaseclient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	clientSlice := make([]Client, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		client := Client{}
		if err := docSnap.DataTo(&client); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		client.Client_id = docSnap.Ref.ID
		clientSlice = append(clientSlice, client)
	}
	return &clientSlice, nil
}

func (db clientFirestore) GetFilteredByUserID(ctx context.Context, userID *string) (*[]Client, *echo.HTTPError) {
	if userID == nil || len(*userID) == 0 {
		return nil, model.ErrFilterNotSet
	}
	query := db.databaseclient.Collection(collection).Query.Where("user_id", "==", *userID)

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

	clientSlice := make([]Client, 0, len(allDocs))
	for _, docSnap := range allDocs {
		clientResponse := Client{}
		if err := docSnap.DataTo(&clientResponse); err != nil {
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}

		clientResponse.Client_id = docSnap.Ref.ID

		clientSlice = append(clientSlice, clientResponse)
	}

	return &clientSlice, nil
}
