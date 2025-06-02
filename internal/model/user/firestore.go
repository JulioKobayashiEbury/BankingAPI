package user

import (
	"errors"
	"net/http"
	"time"

	"BankingAPI/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collection = "users"

type UserFirestore struct {
	Request  *UserRequest
	Response *UserResponse
	Slice    *[]*UserResponse
	AuthUser *AuthenticateUser
	model.Repository
}

func (db *UserFirestore) Create() *model.Erro {
	if db.Request == nil {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()
	entity := map[string]interface{}{
		"name":          db.Request.Name,
		"document":      db.Request.Document,
		"password":      db.Request.Password,
		"register_date": time.Now().Format(model.TimeLayout),
		"status":        true,
	}
	docRef, _, err := clientDB.Collection(collection).Add(*ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	db.Response.User_id = docRef.ID
	return nil
}

func (db *UserFirestore) Delete() *model.Erro {
	if db.Request.User_id == "" {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return model.FailCreatingClient
	}
	defer clientDB.Close()
	docRef := clientDB.Collection(collection).Doc(db.Request.User_id)
	_, err = docRef.Delete(*ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *UserFirestore) Get() *model.Erro {
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docSnapshot, err := clientDB.Collection(collection).Doc(db.Request.User_id).Get(*ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return model.IDnotFound
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + db.Request.User_id)
		return &model.Erro{Err: errors.New("Nil account from snapshot" + (db.Request.User_id)), HttpCode: http.StatusInternalServerError}
	}
	if err := docSnapshot.DataTo(db.Response); err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *UserFirestore) Update() *model.Erro {
	if db.GetUpdateList() == nil {
		log.Error().Msg(model.ResquestNotSet.Err.Error())
		return model.ResquestNotSet
	}
	updates := make([]firestore.Update, 0, 0)
	for key, value := range *db.GetUpdateList() {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docRef := clientDB.Collection((collection)).Doc(db.Request.User_id)

	docSnap, _ := docRef.Get(*ctx)
	if !docSnap.Exists() {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID from collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	_, err = docRef.Update(*ctx, updates)
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db *UserFirestore) GetAll() *model.Erro {
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	iterator := clientDB.Collection(collection).Documents(*ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	userResponseSlice := make([]*UserResponse, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		userResponse := &UserResponse{}
		if err := docSnap.DataTo(&userResponse); err != nil {
			log.Error().Msg(err.Error())
			return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		userResponse.User_id = docSnap.Ref.ID
		// condicional para saber se a transferencia pertence ao account
		userResponseSlice = append(userResponseSlice, userResponse)
	}
	db.Slice = &userResponseSlice
	return nil
}

func (db *UserFirestore) GetAuthInfo() *model.Erro {
	ctx, clientDB, err := model.GetFireStoreClient()
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	defer clientDB.Close()

	docSnapshot, err := clientDB.Collection(collection).Doc(db.Request.User_id).Get(*ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return &model.Erro{Err: errors.New("ID in collection: " + collection + " not found"), HttpCode: http.StatusBadRequest}
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + db.Request.User_id)
		return &model.Erro{Err: errors.New("Nil account from snapshot" + (db.Request.User_id)), HttpCode: http.StatusInternalServerError}
	}
	if err := docSnapshot.DataTo(db.AuthUser); err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}
