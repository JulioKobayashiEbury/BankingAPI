package user

import (
	"BankingAPI/internal/model"
)

/*
Create(interface{}) (*string, *Erro)
Delete(*string) *Erro
Get(id *string) (interface{}, *Erro)
Update(interface{}) *Erro
GetAll() (interface{}, *Erro)
*/

var singleton model.RepositoryInterface

type MockUserRepository struct {
	UserMap *map[string]User
}

func NewMockUserRepository() model.RepositoryInterface {
	if singleton != nil {
		return singleton
	}
	userMap := make(map[string]User)
	singleton = MockUserRepository{
		UserMap: &userMap,
	}
	return singleton
}

func (m MockUserRepository) Create(request interface{}) (interface{}, *model.Erro) {
	userRequest, ok := request.(*User)
	if !ok {
		return nil, model.DataTypeWrong
	}
	/*
		for {
			userID := randomstring.String(10)
			if _, ok := (*m.UserMap)[userID]; !ok {
				userRequest.User_id = userID
				(*m.UserMap)[userID] = *userRequest
				break
			}
		}
	*/
	(*m.UserMap)[userRequest.User_id] = *userRequest
	return userRequest, nil
}

func (m MockUserRepository) Delete(id *string) *model.Erro {
	if _, ok := (*m.UserMap)[*id]; !ok {
		return model.IDnotFound
	}
	delete(*m.UserMap, *id)
	return nil
}

func (m MockUserRepository) Get(id *string) (interface{}, *model.Erro) {
	if user, ok := (*m.UserMap)[*id]; !ok {
		return nil, model.IDnotFound
	} else {
		return &user, nil
	}
}

func (m MockUserRepository) Update(request interface{}) *model.Erro {
	userRequest, ok := request.(*User)
	if !ok {
		return model.DataTypeWrong
	}
	if _, ok := (*m.UserMap)[userRequest.User_id]; !ok {
		return model.IDnotFound
	}
	(*m.UserMap)[userRequest.User_id] = *userRequest
	return nil
}

func (m MockUserRepository) GetAll() (interface{}, *model.Erro) {
	if len(*m.UserMap) == 0 {
		return nil, model.IDnotFound
	}
	users := make([]User, 0, len(*m.UserMap))
	for _, user := range *m.UserMap {
		users = append(users, user)
	}
	return users, nil
}

func (db MockUserRepository) GetFiltered(filters *[]string) (interface{}, *model.Erro) {
	if filters == nil || len(*filters) == 0 {
		return nil, model.FilterNotSet
	}
	userSlice := make([]User, 0, len(*db.UserMap))
	for _, userResponse := range *db.UserMap {
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
