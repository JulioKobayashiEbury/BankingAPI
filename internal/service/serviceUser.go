package service

import (
	"context"
	"time"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"

	"github.com/rs/zerolog/log"
)

type userServiceImpl struct {
	userDatabase   user.UserRepository
	clientDatabase client.ClientRepository
}

func NewUserService(userDB user.UserRepository, clientDB client.ClientRepository) UserService {
	return userServiceImpl{
		userDatabase:   userDB,
		clientDatabase: clientDB,
	}
}

func (service userServiceImpl) Create(ctx context.Context, userRequest *user.User) (*user.User, *model.Erro) {
	if userRequest.Name == "" || userRequest.Document == "" || userRequest.Password == "" {
		log.Warn().Msg("Missing informations on creating user")
		return nil, ErrorMissingCredentials
	}
	userResponse, err := service.userDatabase.Create(ctx, userRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("User created: " + userResponse.User_id)
	return userResponse, nil
}

func (service userServiceImpl) Delete(ctx context.Context, id *string) *model.Erro {
	if err := service.userDatabase.Delete(ctx, id); err != nil {
		return err
	}
	log.Info().Msg("User deleted: " + *id)
	return nil
}

func (service userServiceImpl) Get(ctx context.Context, id *string) (*user.User, *model.Erro) {
	userResponse, err := service.userDatabase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

func (service userServiceImpl) Update(ctx context.Context, userRequest *user.User) (*user.User, *model.Erro) {
	userResponse, err := service.Get(ctx, &(userRequest.User_id))
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

	if err := service.userDatabase.Update(ctx, userResponse); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (user): " + userRequest.User_id)

	return service.Get(ctx, &userRequest.User_id)
}

func (service userServiceImpl) GetAll(ctx context.Context) (*[]user.User, *model.Erro) {
	users, err := service.userDatabase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service userServiceImpl) Report(ctx context.Context, userId *string) (*user.UserReport, *model.Erro) {
	userInfo, err := service.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	clients, err := service.clientDatabase.GetFilteredByID(ctx, userId)
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
