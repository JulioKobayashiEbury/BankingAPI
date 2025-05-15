package controller

import (
	controller "BankingAPI/code/internal/controller/objects"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddClientsEndPoints(server *echo.Echo) {
	server.POST("/clients", ClientPostHandler)
	server.GET("/clients/:client_id", ClientGetHandler)
	server.DELETE("/clients/:client_id", ClientDeleteHandler)
	server.PUT("/clients/:client_id", ClientPutHandler)
	server.PUT("/clients/:client_id/block", ClientPutBlockHandler)
	server.PUT("/clients/:client_id/unblock", ClientPutUnBlockHandler)
}

func ClientPostHandler(c echo.Context) error {
	var clientInfo controller.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientGetHandler(c echo.Context) error {
	var clientInfo controller.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientDeleteHandler(c echo.Context) error {
	var clientInfo controller.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientPutHandler(c echo.Context) error {
	var clientInfo controller.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// clientInfo.UserId = uint32(userID)
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientPutBlockHandler(c echo.Context) error {
	var clientInfo controller.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientPutUnBlockHandler(c echo.Context) error {
	var clientInfo controller.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, controller.ClientRequest{Client_id: clientInfo.Client_id, Status: clientInfo.Status})
}
