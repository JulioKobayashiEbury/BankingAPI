package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collection    = "users"
	cacheDuration = time.Minute * 30
)

var singleton *userFirestore

type userFirestore struct {
	databaseClient *firestore.Client
	cache          map[string]User
	mutex          *sync.Mutex
}

// torna o sistema em statefull
func (db userFirestore) manageDataOnCache(dataID *string) {
	var timeLimit int64 = time.Now().Unix() + int64(cacheDuration)
	for {
		if time.Now().Unix() > timeLimit {
			user, ok := db.cache[*dataID]
			if !ok {
				// log.Info().Msg("user already deleted from cachce")
				return
			}
			ctx := context.Background()
			for {
				if _, err := db.databaseClient.Collection(collection).Doc(*dataID).Set(ctx, &user); err != nil {
					log.Error().Msg("error updating user on DB")
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

func NewUserFireStore(dbClient *firestore.Client) UserRepository {
	if singleton == nil {
		singleton = &userFirestore{
			databaseClient: dbClient,
			cache:          map[string]User{},
			mutex:          &sync.Mutex{},
		}
	}
	return singleton
}

func (db userFirestore) Create(ctx context.Context, request *User) (*User, *echo.HTTPError) {
	db.mutex.Lock()
	db.cache[request.User_id] = *request
	go db.manageDataOnCache(&request.User_id)
	db.mutex.Unlock()

	entity := map[string]interface{}{
		"name":          request.Name,
		"document":      request.Document,
		"password":      request.Password,
		"register_date": time.Now().Format(model.TimeLayout),
		"status":        model.ValidStatus[0],
	}
	docRef, _, err := db.databaseClient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return db.Get(ctx, &docRef.ID)
}

func (db userFirestore) Delete(ctx context.Context, id *string) *echo.HTTPError {
	ch := make(chan *echo.HTTPError)

	go func(ch chan *echo.HTTPError) {
		docRef := db.databaseClient.Collection(collection).Doc(*id)

		if _, err := docRef.Delete(ctx); err != nil {
			log.Error().Msg(err.Error())
			ch <- &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
	}(ch)

	db.mutex.Lock()
	delete(db.cache, *id)
	db.mutex.Unlock()
	err := <-ch
	close(ch)
	return err
}

func (db userFirestore) Get(ctx context.Context, id *string) (*User, *echo.HTTPError) {
	db.mutex.Lock()
	if response, ok := db.cache[*id]; ok {
		log.Info().Msg("retrieved user from cache")
		db.mutex.Unlock()
		return &response, nil
	}
	db.mutex.Unlock()

	docSnapshot, err := db.databaseClient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, model.ErrIDnotFound
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &echo.HTTPError{Internal: errors.New("Nil account from snapshot" + (*id)), Code: http.StatusInternalServerError, Message: fmt.Sprint("Nil account from snapshot" + (*id))}
	}
	userResponse := User{}
	if err := docSnapshot.DataTo(&userResponse); err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	userResponse.User_id = docSnapshot.Ref.ID

	db.mutex.Lock()
	db.cache[userResponse.User_id] = userResponse
	go db.manageDataOnCache(&userResponse.User_id)
	db.mutex.Unlock()

	return &userResponse, nil
}

func (db userFirestore) Update(ctx context.Context, request *User) *echo.HTTPError {
	ch := make(chan *echo.HTTPError)

	go func(errChan chan *echo.HTTPError) {
		entity := map[string]interface{}{
			"name":          request.Name,
			"document":      request.Document,
			"password":      request.Password,
			"register_date": request.Register_date,
		}

		docRef := db.databaseClient.Collection(collection).Doc(request.User_id)
		if _, err := docRef.Set(ctx, entity); err != nil {
			db.mutex.Lock()
			delete(db.cache, request.User_id)
			db.mutex.Unlock()
			log.Error().Msg(err.Error())
			errChan <- &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}

		log.Info().Msg("user: " + request.User_id + "has been updated")

		errChan <- nil
	}(ch)

	db.mutex.Lock()
	if _, ok := db.cache[request.User_id]; ok {
		db.cache[request.User_id] = *request
	}
	db.mutex.Unlock()

	err := <-ch
	close(ch)
	return err
}

func (db userFirestore) GetAll(ctx context.Context) (*[]User, *echo.HTTPError) {
	iterator := db.databaseClient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	userResponseSlice := make([]User, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		userResponse := User{}
		if err := docSnap.DataTo(&userResponse); err != nil {
			log.Error().Msg(err.Error())
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}
		userResponse.User_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		userResponseSlice = append(userResponseSlice, userResponse)
	}
	return &userResponseSlice, nil
}

// "key,==.value"

func (db userFirestore) GetFilteredByID(ctx context.Context, filters *string) (*[]User, *echo.HTTPError) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.ErrFilterNotSet
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
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}

	userSlice := make([]User, 0, len(allDocs))
	for _, docSnap := range allDocs {
		userResponse := User{}
		if err := docSnap.DataTo(&userResponse); err != nil {
			return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
		}

		userResponse.User_id = docSnap.Ref.ID

		userSlice = append(userSlice, userResponse)
	}
	return &userSlice, nil
}
