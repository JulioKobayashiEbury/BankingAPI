package controller

import (
	"net/http"
	"strconv"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/service"

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
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	clientResponse, err := service.CreateClient(&clientInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientGetHandler(c echo.Context) error {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	clientInfo, err := service.GetClient(uint32(clientID))
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func ClientDeleteHandler(c echo.Context) error {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.ClientDelete(uint32(clientID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func ClientPutHandler(c echo.Context) error {
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	clientResponse, err := service.UpdateClient(&clientInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientPutBlockHandler(c echo.Context) error {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.ClientBlock(uint32(clientID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Blocked"})
}

func ClientPutUnBlockHandler(c echo.Context) error {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.ClientUnBlock(uint32(clientID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Unblocked"})
}
