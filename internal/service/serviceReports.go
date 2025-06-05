package service

import (
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/user"
)

type ReportService interface {
	GenerateReportByAccount(accountID *string) (*account.AccountReport, *model.Erro)
	GenerateReportByClient(clientID *string) (*client.ClientReport, *model.Erro)
	GenerateReportByUser(userId *string) (*user.UserReport, *model.Erro)
}

type reportImpl struct {
	getService    ServiceGet
	getAllService ServiceGetAll
}

func NewReportService(toGetService ServiceGet, toGetAllService ServiceGetAll) ReportService {
	return reportImpl{
		getService:    toGetService,
		getAllService: toGetAllService,
	}
}
func (report reportImpl) GenerateReportByClient(clientID *string) (*client.ClientReport, *model.Erro) {
	clientInfo, err := report.getService.Client(*clientID)
	if err != nil {
		return nil, err
	}
	accounts, err := report.getAllService.GetAccountsByClientID(clientID)
	if err != nil {
		return nil, err
	}
	return &client.ClientReport{
		Client_id:     clientInfo.Client_id,
		User_id:       clientInfo.User_id,
		Name:          clientInfo.Name,
		Document:      clientInfo.Document,
		Register_date: clientInfo.Register_date,
		Status:        clientInfo.Status,
		Accounts:      (*accounts),
		Report_date:   time.Now().Format(timeLayout),
	}, nil
}

func (report reportImpl) GenerateReportByUser(userId *string) (*user.UserReport, *model.Erro) {
	userInfo, err := report.getService.User(*userId)
	if err != nil {
		return nil, err
	}
	clients, err := report.getAllService.GetClientsByUserID(userId)
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
