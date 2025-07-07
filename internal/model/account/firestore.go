package account

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	model "BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collection = "accounts"

var (
	singleton     *accountFirestore
	cacheDuration = time.Minute * 1
)

type accountFirestore struct {
	databaseClient *firestore.Client
	cache          map[string]Account
	mutex          *sync.Mutex
}

func (db accountFirestore) manageDataOnCache(dataID *string) {
	var timeLimit int64 = time.Now().Unix() + int64(cacheDuration)
	for {
		if time.Now().Unix() > timeLimit {
			account, ok := db.cache[*dataID]
			if !ok {
				// log.Info().Msg("user already deleted from cachce")
				return
			}
			ctx := context.Background()
			for {
				if _, err := db.databaseClient.Collection(collection).Doc(*dataID).Set(ctx, &account); err != nil {
					log.Error().Msg("error updating account on DB")
					time.Sleep(time.Duration(1))
					continue
				}
				break
			}
			ctx.Done()
			delete(db.cache, *dataID)
			// log.Info().Msg("stopped routine for data on cache")
			return
		}
	}
}

func NewAccountFirestore(dbClient *firestore.Client) AccountRepository {
	return accountFirestore{
		databaseClient: dbClient,
		cache:          make(map[string]Account),
		mutex:          &sync.Mutex{},
	}
}

func (db accountFirestore) Create(ctx context.Context, request *Account) (*Account, *echo.HTTPError) {
	entity := map[string]interface{}{
		"client_id":     request.Client_id,
		"user_id":       request.User_id,
		"agency_id":     request.Agency_id,
		"balance":       0.0,
		"register_date": time.Now().Format(model.TimeLayout),
		"status":        model.ValidStatus[0],
	}

	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	request.Account_id = docRef.ID

	db.mutex.Lock()
	db.cache[docRef.ID] = *request
	go db.manageDataOnCache(&docRef.ID)
	db.mutex.Unlock()

	return db.Get(ctx, &docRef.ID)
}

func (db accountFirestore) Delete(ctx context.Context, id *string) *echo.HTTPError {
	ch := make(chan *echo.HTTPError, 1)
	go func(ch chan *echo.HTTPError) {
		docRef := db.databaseClient.Collection(collection).Doc(*id)

		if _, err := docRef.Delete(ctx); err != nil {
			log.Error().Msg(err.Error())
			ch <- &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}

		ch <- nil
	}(ch)
	db.mutex.Lock()
	delete(db.cache, *id)
	db.mutex.Unlock()

	err := <-ch
	close(ch)
	return err
}

func (db accountFirestore) Get(ctx context.Context, id *string) (*Account, *echo.HTTPError) {
	db.mutex.Lock()
	if response, ok := db.cache[*id]; ok {
		log.Info().Msg("retrievced account from cache")
		db.mutex.Unlock()
		return &response, nil
	}
	db.mutex.Unlock()

	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, &echo.HTTPError{Internal: errors.New("ID in collection: " + collection + " not found"), Code: http.StatusBadRequest, Message: fmt.Sprintf("ID from collection: " + collection + " not found")}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &echo.HTTPError{Internal: errors.New("Nil account from snapshot" + *id), Code: http.StatusInternalServerError, Message: fmt.Sprint("Nil account from snapshot" + *id)}
	}
	accountResponse := Account{}
	if err := docSnapshot.DataTo(&accountResponse); err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	accountResponse.Account_id = docSnapshot.Ref.ID

	db.mutex.Lock()
	db.cache[accountResponse.Account_id] = accountResponse
	go db.manageDataOnCache(&accountResponse.Account_id)
	db.mutex.Unlock()

	return &accountResponse, nil
}

func (db accountFirestore) Update(ctx context.Context, request *Account) *echo.HTTPError {
	ch := make(chan *echo.HTTPError)
	go func(ch chan *echo.HTTPError) {
		entity := map[string]interface{}{
			"client_id":     request.Client_id,
			"user_id":       request.User_id,
			"agency_id":     request.Agency_id,
			"balance":       request.Balance,
			"register_date": request.Register_date,
			"status":        request.Status,
		}
		docRef := db.databaseClient.Collection(collection).Doc(request.Account_id)

		if _, err := docRef.Set(ctx, entity); err != nil {
			db.mutex.Lock()
			delete(db.cache, request.Account_id)
			db.mutex.Unlock()
			log.Error().Msg(err.Error())
			ch <- &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		log.Info().Msg("Account: " + request.Account_id + " has been updated")
		ch <- nil
	}(ch)

	db.mutex.Lock()
	if _, ok := db.cache[request.Account_id]; ok {
		db.cache[request.Account_id] = *request
	}
	db.mutex.Unlock()
	err := <-ch
	close(ch)
	return err
}

func (db accountFirestore) GetAll(ctx context.Context) (*[]Account, *echo.HTTPError) {
	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	accountResponseSlice := make([]Account, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		accountResponse := Account{}
		if err := docSnap.DataTo(&accountResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		accountResponse.Account_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		accountResponseSlice = append(accountResponseSlice, accountResponse)
	}
	return &accountResponseSlice, nil
}

func (db accountFirestore) GetFilteredByClientID(ctx context.Context, clientID *string) (*[]Account, *echo.HTTPError) {
	//"user_id,==,1"
	if clientID == nil || len(*clientID) == 0 {
		return nil, model.ErrFilterNotSet
	}

	query := db.databaseClient.Collection(collection).Query.Where("client_id", "==", *clientID)

	/* for _, filter := range *filters {
		tokens := model.TokenizeFilters(&filter)
		if len(*tokens) != 3 {
			log.Error().Msg("Invalid filter format: " + filter)
			return nil, model.InvalidFilterFormat
		}
		query = query.Where((*tokens)[0], (*tokens)[1], (*tokens)[2])
	}
	*/

	allDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError}
	}
	accountResponseSlice := make([]Account, 0, len(allDocs))
	for _, docSnap := range allDocs {
		accountResponse := Account{}
		if err := docSnap.DataTo(&accountResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		accountResponse.Account_id = docSnap.Ref.ID

		accountResponseSlice = append(accountResponseSlice, accountResponse)

	}

	return &accountResponseSlice, nil
}
