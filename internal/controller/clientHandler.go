package controller

import (
	"net/http"

	"BankingAPI/internal/model/client"
	"BankingAPI/internal/service"

	model "BankingAPI/internal/model"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type ClientHandler interface {
	ClientPostHandler(c echo.Context) error
	ClientGetHandler(c echo.Context) error
	ClientDeleteHandler(c echo.Context) error
	ClientPutHandler(c echo.Context) error
	ClientGetReportHandler(c echo.Context) error
}

type clientHandlerImpl struct {
	clientService service.ClientService
}

func NewClientHandler(clientService service.ClientService) ClientHandler {
	return clientHandlerImpl{
		clientService: clientService,
	}
}

func AddClientsEndPoints(server *echo.Echo, h ClientHandler) {
	server.POST("/clients", h.ClientPostHandler)
	server.GET("/clients/:client_id", h.ClientGetHandler)
	server.GET("/clients/:client_id/report", h.ClientGetReportHandler)
	server.DELETE("/clients/:client_id", h.ClientDeleteHandler)
	server.PUT("/clients/:client_id", h.ClientPutHandler)
}

func (h clientHandlerImpl) ClientPostHandler(c echo.Context) error {
	var clientInfo client.Client
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(clientInfo.Document) != documentLenghtForClient || len(clientInfo.Name) > maxNameLenght || clientInfo.User_id == "" {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	Client, err := h.clientService.Create(c.Request().Context(), &clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*Client))
}

func (h clientHandlerImpl) ClientGetHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	clientInfo, err := h.clientService.Get(c.Request().Context(), &clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func (h clientHandlerImpl) ClientDeleteHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if err := h.clientService.Delete(c.Request().Context(), &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func (h clientHandlerImpl) ClientPutHandler(c echo.Context) error {
	clientID := c.Param("client_id")

	var clientInfo client.Client
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	clientInfo.Client_id = clientID
	Client, err := h.clientService.Update(c.Request().Context(), &clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*Client))
}

func (h clientHandlerImpl) ClientGetReportHandler(c echo.Context) error {
	clientID := c.Param("client_id")

	clientReport, err := h.clientService.Report(c.Request().Context(), &clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*clientReport))
}
