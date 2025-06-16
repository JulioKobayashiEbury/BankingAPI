package controller

import (
	"errors"
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
	userID, err := h.authorizationForClientEndpoints(&c, nil)
	if err != nil {
		c.JSON(err.HttpCode, err.Err.Error())
	}

	var clientInfo client.Client
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if *userID != clientInfo.User_id {
		log.Warn().Msg("User ID does not match with clients User ID")
		return c.JSON(http.StatusForbidden, model.StandartResponse{Message: "User ID does not match with clients User ID"})
	}

	if len(clientInfo.Document) != documentLenghtForClient || len(clientInfo.Name) > maxNameLenght {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	Client, err := h.clientService.Create(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*Client))
}

func (h clientHandlerImpl) ClientGetHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := h.authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	clientInfo, err := h.clientService.Get(&clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func (h clientHandlerImpl) ClientDeleteHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := h.authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.clientService.Delete(&clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func (h clientHandlerImpl) ClientPutHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := h.authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	var clientInfo client.Client
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	clientInfo.Client_id = clientID
	Client, err := h.clientService.Update(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*Client))
}

func (h clientHandlerImpl) ClientGetReportHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := h.authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	clientReport, err := h.clientService.Report(&clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*clientReport))
}

func (h clientHandlerImpl) authorizationForClientEndpoints(c *echo.Context, clientID *string) (*string, *model.Erro) {
	authorizationHeader := (*c).Request().Header.Get((echo.HeaderAuthorization))

	claims, err := service.Authorize(&authorizationHeader)
	if err != nil {
		return nil, err
	}

	if clientID == nil {
		return &claims.Id, nil
	}

	client, err := h.clientService.Get(clientID)
	if err != nil {
		return nil, err
	}

	if client.User_id != claims.Id {
		log.Error().Msg("User ID does not match with accounts User ID")
		return nil, &model.Erro{Err: errors.New("no match for user id"), HttpCode: http.StatusForbidden}
	}

	return nil, nil
}
