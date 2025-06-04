package controller

import (
	"net/http"

	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/model/withdrawal"
	"BankingAPI/internal/service"

	model "BankingAPI/internal/model"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddClientsEndPoints(server *echo.Echo) {
	server.POST("/clients", ClientPostHandler)
	server.GET("/clients/:client_id", ClientGetHandler)
	server.GET("/clients/:client_id/report", ClientGetReportHandler)
	server.DELETE("/clients/:client_id", ClientDeleteHandler)
	server.PUT("/clients/:client_id", ClientPutHandler)
	server.PUT("/clients/:client_id/block", ClientPutBlockHandler)
	server.PUT("/clients/:client_id/unblock", ClientPutUnBlockHandler)
}

func ClientPostHandler(c echo.Context) error {
	var clientInfo client.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(clientInfo.Document) != documentLenghtIdeal || len(clientInfo.Name) > maxNameLenght {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	userDatabase := user.NewUserFireStore(DatabaseClient)

	serviceGet := service.NewGetService(nil, clientDatabase, userDatabase)
	serviceCreate := service.NewCreateService(nil, clientDatabase, userDatabase, serviceGet)
	clientResponse, err := serviceCreate.CreateClient(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientGetHandler(c echo.Context) error {
	_, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientID := c.Param("client_id")

	clientDatabase := client.NewClientFirestore(DatabaseClient)
	serviceGet := service.NewGetService(nil, clientDatabase, nil)
	clientInfo, err := serviceGet.Client(clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func ClientDeleteHandler(c echo.Context) error {
	_, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientID := c.Param("client_id")
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	serviceDelete := service.NewDeleteService(nil, clientDatabase, nil)
	if err := serviceDelete.ClientDelete(clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func ClientPutHandler(c echo.Context) error {
	_, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var clientInfo client.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	clientInfo.Client_id = c.Param("client_id")
	userDatabase := user.NewUserFireStore(DatabaseClient)
	clientDatabase := client.NewClientFirestore(DatabaseClient)

	serviceGet := service.NewGetService(nil, clientDatabase, userDatabase)
	serviceUpdate := service.NewUpdateService(nil, clientDatabase, userDatabase, serviceGet)

	clientResponse, err := serviceUpdate.UpdateClient(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientPutBlockHandler(c echo.Context) error {
	_, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientID := c.Param("client_id")
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	serviceGet := service.NewGetService(nil, clientDatabase, nil)
	serviceStatus := service.NewStatusService(nil, clientDatabase, nil, serviceGet)
	if err := serviceStatus.ClientBlock(clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Blocked"})
}

func ClientPutUnBlockHandler(c echo.Context) error {
	_, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientID := c.Param("client_id")
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	serviceGet := service.NewGetService(nil, clientDatabase, nil)
	serviceStatus := service.NewStatusService(nil, clientDatabase, nil, serviceGet)
	if err := serviceStatus.ClientUnBlock(clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Unblocked"})
}

func ClientGetReportHandler(c echo.Context) error {
	_, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientID := c.Param("client_id")

	autodebitDatabase := automaticdebit.NewAutoDebitFirestore(DatabaseClient)
	withdrawalDatabase := withdrawal.NewWithdrawalFirestore(DatabaseClient)
	depositDatabase := deposit.NewDepositFirestore(DatabaseClient)
	transferDatabase := transfer.NewTransferFirestore(DatabaseClient)
	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	userDatabase := user.NewUserFireStore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, clientDatabase, userDatabase)
	serviceGetAll := service.NewGetAllService(autodebitDatabase, withdrawalDatabase, depositDatabase, transferDatabase, accountDatabase, clientDatabase)

	serviceReport := service.NewReportService(serviceGet, serviceGetAll)

	clientReport, err := serviceReport.GenerateReportByClient(&clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*clientReport))
}
