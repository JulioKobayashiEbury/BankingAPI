package controller

import (
	"net/http"

	"BankingAPI/internal/model/client"
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
	clientResponse, err := service.CreateClient(&clientInfo)
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
	clientInfo, err := service.Client(clientID)
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
	if err := service.ClientDelete(clientID); err != nil {
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

	clientResponse, err := service.UpdateClient(&clientInfo)
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
	if err := service.ClientBlock(clientID); err != nil {
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
	if err := service.ClientUnBlock(clientID); err != nil {
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
	clientReport, err := service.GenerateReportByClient(&clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*clientReport))
}
