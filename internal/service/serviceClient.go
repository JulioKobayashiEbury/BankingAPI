package service

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/model/client"
)

type userServiceImpl struct {
	clientDatabase model.RepositoryInterface
	userService    UserService
}

func NewUserService(cd model.RepositoryInterface, us UserService) ClientService {
	return userServiceImpl{
		clientDatabase: cd,
		userService:    us,
	}
}

func (service userServiceImpl) Create(*client.Client) (*string, *model.Erro) {

}

func (service userServiceImpl) Delete(*string) *model.Erro
func (service userServiceImpl) Get(*string) (*client.Client, *model.Erro)
func (service userServiceImpl) Update(*client.Client) *model.Erro
func (service userServiceImpl) GetAll() ([]*client.Client, *model.Erro)
func (service userServiceImpl) Status(*string, bool) *model.Erro
func (service userServiceImpl) Report(*string) (*client.ClientReport, *model.Erro)
func (service userServiceImpl) GetClientsByUserID(userID *string) (*[]client.Client, *model.Erro)
