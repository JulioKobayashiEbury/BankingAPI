package user

import (
	"errors"
	"net/http"
	"sync"

	"BankingAPI/internal/model"

	"github.com/rs/zerolog/log"
)

/*
Create(interface{}) (*string, *Erro)
Delete(*string) *Erro
Get(id *string) (interface{}, *Erro)
Update(interface{}) *Erro
GetAll() (interface{}, *Erro)
*/

var (
	once      sync.Once
	singleton model.RepositoryInterface
)

type MockUserRepository struct {
	idList   *[]string
	usersMap *map[string]User
}

func NewMockUserRepository() model.RepositoryInterface {
	once.Do(func() {
		mapUsers := make(map[string]User)
		listId := make([]string, 0)
		singleton = MockUserRepository{
			idList:   &listId,
			usersMap: &mapUsers,
		}
	})
	return singleton
}

func (m MockUserRepository) Create(request interface{}) (interface{}, *model.Erro) {
	user, ok := request.(*User)
	if !ok {
		return nil, model.DataTypeWrong
	}
	if _, ok := (*m.usersMap)[user.User_id]; ok {
		return nil, &model.Erro{Err: errors.New("id already exists"), HttpCode: http.StatusBadRequest}
	}

	(*m.idList) = append((*m.idList), user.User_id)
	(*m.usersMap)[user.User_id] = *user

	return request, nil
}

func (m MockUserRepository) Delete(id *string) *model.Erro {
	if id == nil || *id == "" {
		return model.IDnotFound
	}
	if _, ok := (*m.usersMap)[*id]; !ok {
		log.Debug().Msg("No user in usermap")
		return model.IDnotFound
	}
	delete(*m.usersMap, *id)
	for index, userId := range *m.idList {
		if userId == *id {
			deleteIndex := index
			// Remove the user ID from the list
			for i := deleteIndex; i < len(*m.idList)-1; i++ {
				if i+1 < len(*m.idList) {
					(*m.idList)[i] = (*m.idList)[i+1]
				}
			}
			// Move the last element to the current index
			// Resize the slice to remove the last element
			break
		}
	}
	return nil
}

func (m MockUserRepository) Get(id *string) (interface{}, *model.Erro) {
	if id == nil || *id == "" {
		return nil, model.ResquestNotSet
	}
	for _, userId := range *m.idList {
		if userId == *id {
			if user, ok := (*m.usersMap)[*id]; ok {
				return &user, nil
			}
		}
	}
	return nil, model.IDnotFound
}

func (m MockUserRepository) Update(request interface{}) *model.Erro {
	user, ok := request.(*User)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*m.usersMap)[user.User_id]; !ok {
		return model.IDnotFound
	}
	(*m.usersMap)[user.User_id] = *user
	return nil
}

func (m MockUserRepository) GetAll() (interface{}, *model.Erro) {
	if len(*m.usersMap) == 0 {
		return nil, model.IDnotFound
	}
	users := make([]User, 0, len(*m.usersMap))
	for _, userId := range *m.idList {
		user, ok := (*m.usersMap)[userId]
		if !ok {
			return nil, model.IDnotFound
		}
		users = append(users, user)
	}
	return &users, nil
}

func (db MockUserRepository) GetFiltered(filters *[]string) (interface{}, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	userSlice := make([]User, 0, len(*db.usersMap))
	for _, userId := range *db.idList {
		userResponse := (*db.usersMap)[userId]
		match := true
		for _, filter := range *filters {
			token := model.TokenizeFilters(&filter)
			if len(*token) != 3 {
				return nil, model.InvalidFilterFormat
			}
			field := (*token)[0]
			operator := (*token)[1]
			value := (*token)[2]

			switch field {
			case "register_date":
				if operator == ">=" && userResponse.Register_date < value {
					match = false
				}
				if operator == "<=" && userResponse.Register_date > value {
					match = false
				}
			// Add more fields as necessary
			default:
				match = false
			}
		}
		if match {
			userSlice = append(userSlice, userResponse)
		}
	}
	return &userSlice, nil
}
