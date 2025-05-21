package controller

import (
	"net/http"

	"BankingAPI/internal/service"

	model "BankingAPI/internal/model/types"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddClientsEndPoints(server *echo.Echo) {
	server.POST("/clients", ClientPostHandler)
	server.POST("/clients/auth", ClientAuthHandler)
	server.GET("/clients/:client_id", ClientGetHandler)
	server.DELETE("/clients/:client_id", ClientDeleteHandler)
	server.PUT("/clients/:client_id", ClientPutHandler)
	server.PUT("/clients/:client_id/block", ClientPutBlockHandler)
	server.PUT("/clients/:client_id/unblock", ClientPutUnBlockHandler)
}

func ClientPostHandler(c echo.Context) error {
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	clientResponse, err := service.CreateClient(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientAuthHandler(c echo.Context) error {
	// work here
	return nil
}

func ClientGetHandler(c echo.Context) error {
	clientID := c.Param("client_id")

	clientInfo, err := service.Client(clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func ClientDeleteHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if err := service.ClientDelete(clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func ClientPutHandler(c echo.Context) error {
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	clientInfo.Client_id = c.Param("client_id")

	clientResponse, err := service.UpdateClient(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientPutBlockHandler(c echo.Context) error {
	clientID := c.Param("client_id")

	if err := service.ClientBlock(string(clientID)); err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Blocked"})
}

func ClientPutUnBlockHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if err := service.ClientUnBlock(clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Unblocked"})
}
