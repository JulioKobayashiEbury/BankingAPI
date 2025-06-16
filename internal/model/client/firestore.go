package client

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

const collection = "clients"

type clientFirestore struct {
	databaseclient *firestore.Client
}

func NewClientFirestore(dbclient *firestore.Client) model.RepositoryInterface {
	return clientFirestore{
		databaseclient: dbclient,
	}
}

func (db clientFirestore) Create(request interface{}) (interface{}, *model.Erro) {
	client, ok := request.(*Client)
	if !ok {
		return nil, model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"user_id":       client.User_id,
		"name":          client.Name,
		"document":      client.Document,
		"status":        model.ValidStatus[0],
		"register_date": time.Now().Format(model.TimeLayout),
	}
	docRef, _, err := db.databaseclient.Collection(collection).Add(ctx, entity)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return db.Get(&docRef.ID)
}

func (db clientFirestore) Delete(id *string) *model.Erro {
	ctx := context.Background()
	defer ctx.Done()

	docRef := db.databaseclient.Collection(collection).Doc(*id)
	if _, err := docRef.Delete(ctx); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return nil
}

func (db clientFirestore) Get(id *string) (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	docSnapshot, err := db.databaseclient.Collection(collection).Doc(*id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		log.Warn().Msg("ID from collection: " + collection + " not found")
		return nil, model.IDnotFound
	}
	if docSnapshot == nil {
		log.Error().Msg("Nil account from snapshot" + *id)
		return nil, &model.Erro{Err: errors.New("Nil account from snapshot" + (*id)), HttpCode: http.StatusInternalServerError}
	}
	client := Client{}
	if err := docSnapshot.DataTo(&client); err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	client.Client_id = docSnapshot.Ref.ID
	return &client, nil
}

func (db clientFirestore) Update(request interface{}) *model.Erro {
	client, ok := request.(*Client)
	if !ok {
		return model.DataTypeWrong
	}
	ctx := context.Background()
	defer ctx.Done()

	entity := map[string]interface{}{
		"user_id":       client.User_id,
		"name":          client.Name,
		"document":      client.Document,
		"status":        client.Status,
		"register_date": client.Register_date,
	}
	docRef := db.databaseclient.Collection(collection).Doc(client.Client_id)

	if _, err := docRef.Set(ctx, entity); err != nil {
		log.Error().Msg(err.Error())
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	log.Info().Msg("Account: " + client.Client_id + " has been updated")

	return nil
}

func (db clientFirestore) GetAll() (interface{}, *model.Erro) {
	ctx := context.Background()
	defer ctx.Done()

	iterator := db.databaseclient.Collection(collection).Documents(ctx)

	docSnapshots, err := iterator.GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	clientSlice := make([]Client, 0, len(docSnapshots))
	for index := 0; index < len(docSnapshots); index++ {
		docSnap := docSnapshots[index]
		client := Client{}
		if err := docSnap.DataTo(&client); err != nil {
			log.Error().Msg(err.Error())
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}
		client.Client_id = docSnap.Ref.ID
		clientSlice = append(clientSlice, client)
	}
	return &clientSlice, nil
}

func (db clientFirestore) GetFiltered(filters *[]string) (interface{}, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}

	ctx := context.Background()
	defer ctx.Done()

	query := db.databaseclient.Collection(collection).Query

	for _, filter := range *filters {
		token := model.TokenizeFilters(&filter)
		if len(*token) != 3 {
			return nil, model.InvalidFilterFormat
		}

		query = query.Where((*token)[0], (*token)[1], (*token)[2])

	}

	allDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}

	clientSlice := make([]Client, 0, len(allDocs))
	for _, docSnap := range allDocs {
		clientResponse := Client{}
		if err := docSnap.DataTo(&clientResponse); err != nil {
			return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
		}

		clientResponse.Client_id = docSnap.Ref.ID

		clientSlice = append(clientSlice, clientResponse)
	}

	return &clientSlice, nil
}
