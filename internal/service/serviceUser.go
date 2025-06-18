package service

import (
	"time"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

type userServiceImpl struct {
	userDatabase   model.RepositoryInterface
	clientDatabase model.RepositoryInterface
}

func NewUserService(userDB model.RepositoryInterface, clientDB model.RepositoryInterface) UserService {
	return userServiceImpl{
		userDatabase:   userDB,
		clientDatabase: clientDB,
	}
}

func (service userServiceImpl) Create(userRequest *user.User) (*user.User, *model.Erro) {
	if userRequest.Name == "" || userRequest.Document == "" || userRequest.Password == "" {
		log.Warn().Msg("Missing informations on creating user")
		return nil, ErrorMissingCredentials
	}
	obj, err := service.userDatabase.Create(userRequest)
	if err != nil {
		return nil, err
	}
	userResponse, ok := obj.(*user.User)
	if !ok {
		return nil, model.DataTypeWrong
	}

	log.Info().Msg("User created: " + userResponse.User_id)
	return userResponse, nil
}

func (service userServiceImpl) Delete(id *string) *model.Erro {
	if err := service.userDatabase.Delete(id); err != nil {
		return err
	}
	log.Info().Msg("User deleted: " + *id)
	return nil
}

func (service userServiceImpl) Get(id *string) (*user.User, *model.Erro) {
	obj, err := service.userDatabase.Get(id)
	if err != nil {
		return nil, err
	}
	userResponse, ok := obj.(*user.User)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return userResponse, nil
}

func (service userServiceImpl) Update(userRequest *user.User) (*user.User, *model.Erro) {
	userResponse, err := service.Get(&(userRequest.User_id))
	if err != nil {
		return nil, err
	}

	if userRequest.Name != "" {
		userResponse.Name = userRequest.Name
	}
	if userRequest.Document != "" {
		userResponse.Document = userRequest.Document
	}
	if userRequest.Password != "" {
		userResponse.Password = userRequest.Password
	}

	// monta struct de updat

	if err := service.userDatabase.Update(userResponse); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (user): " + userRequest.User_id)

	return service.Get(&userRequest.User_id)
}

func (service userServiceImpl) GetAll() (*[]user.User, *model.Erro) {
	obj, err := service.userDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	users, ok := obj.(*[]user.User)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return users, nil
}

func (service userServiceImpl) Report(id *string) (*user.UserReport, *model.Erro) {
	userInfo, err := service.Get(id)
	if err != nil {
		return nil, err
	}
	clients, err := service.clientDatabase.GetFiltered(&[]string{"user_id,==," + *id})
	if err != nil {
		return nil, err
	}
	return &user.UserReport{
		User_id:       userInfo.User_id,
		Name:          userInfo.Name,
		Document:      userInfo.Document,
		Register_date: userInfo.Register_date,
		Clients:       clients,
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}
