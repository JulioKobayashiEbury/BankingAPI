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

var singleton *model.RepositoryInterface

type MockUserRepository struct {
	UserMap *map[string]User
}

func NewMockUserRepository() model.RepositoryInterface {
	if singleton != nil {
		return *singleton
	}
	userMap := make(map[string]User)
	*singleton = MockUserRepository{
		UserMap: &userMap,
	}
	return *singleton
}

func (m MockUserRepository) Create(request interface{}) (*string, *model.Erro) {
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
	return &userRequest.User_id, nil
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
