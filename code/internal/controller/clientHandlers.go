package controller

import (
	"net/http"
	"strconv"

	"BankingAPI/code/internal/domain"

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
	var clientInfo domain.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	userID, err := strconv.ParseUint(c.FormValue("user_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	clientInfo.UserId = uint32(userID)
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientGetHandler(c echo.Context) error {
	var clientInfo domain.ClientRequest
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	clientInfo.ClientId = uint32(clientID)

	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientDeleteHandler(c echo.Context) error {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, clientID)
}

func ClientPutHandler(c echo.Context) error {
	var clientInfo domain.ClientRequest
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	userID, err := strconv.ParseUint(c.FormValue("user_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err = c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	clientInfo.ClientId = uint32(clientID)
	clientInfo.UserId = uint32(userID)
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientPutBlockHandler(c echo.Context) error {
	var clientInfo domain.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}

func ClientPutUnBlockHandler(c echo.Context) error {
	var clientInfo domain.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, clientInfo)
}
