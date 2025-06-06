package service

import (
	"errors"
	"net/http"
	"time"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

type userServiceImpl struct {
	userDatabase       model.RepositoryInterface
	getFilteredService GetFilteredService
}

func NewUserService(userRepo model.RepositoryInterface, getFilteredServe GetFilteredService) UserService {
	return userServiceImpl{
		userDatabase:       userRepo,
		getFilteredService: getFilteredServe,
	}
}

func (service userServiceImpl) Create(userRequest *user.User) (*string, *model.Erro) {
	if userRequest.Name == "" || userRequest.Document == "" || userRequest.Password == "" {
		log.Warn().Msg("Missing credentials on creating user")
		return nil, &model.Erro{Err: errors.New("Missing credentials"), HttpCode: http.StatusBadRequest}
	}
	userID, err := service.userDatabase.Create(userRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("User created: " + *userID)
	return userID, nil
}

func (service userServiceImpl) Delete(id *string) *model.Erro {
	if err := service.userDatabase.Delete(id); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + *id)
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

func (service userServiceImpl) Status(id *string, status bool) *model.Erro {
	user, err := service.Get(id)
	if err != nil {
		return err
	}
	user.Status = status
	if err := service.userDatabase.Update(user); err != nil {
		return err
	}
	return nil
}

func (service userServiceImpl) Report(id *string) (*user.UserReport, *model.Erro) {
	userInfo, err := service.Get(id)
	if err != nil {
		return nil, err
	}
	clients, err := service.getFilteredService.GetClientsByUserID(id)
	if err != nil {
		return nil, err
	}
	return &user.UserReport{
		User_id:       userInfo.User_id,
		Name:          userInfo.Name,
		Document:      userInfo.Document,
		Register_date: userInfo.Register_date,
		Status:        userInfo.Status,
		Clients:       *clients,
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}
