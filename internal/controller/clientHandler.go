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
	userID, err := authorizationForClientEndpoints(&c, nil)
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

	if len(clientInfo.Document) != documentLenghtIdeal || len(clientInfo.Name) > maxNameLenght {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	Client, err := Services.ClientService.Create(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*Client))
}

func ClientGetHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	clientInfo, err := Services.ClientService.Get(&clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func ClientDeleteHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := Services.ClientService.Delete(&clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func ClientPutHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	var clientInfo client.Client
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	clientInfo.Client_id = clientID
	Client, err := Services.ClientService.Update(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*Client))
}

func ClientPutBlockHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := Services.ClientService.Status(&clientID, false); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Blocked"})
}

func ClientPutUnBlockHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := Services.ClientService.Status(&clientID, true); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Unblocked"})
}

func ClientGetReportHandler(c echo.Context) error {
	clientID := c.Param("client_id")
	if _, err := authorizationForClientEndpoints(&c, &clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	clientReport, err := Services.ClientService.Report(&clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*clientReport))
}

func authorizationForClientEndpoints(c *echo.Context, clientID *string) (*string, *model.Erro) {
	authorizationHeader := (*c).Request().Header.Get((echo.HeaderAuthorization))

	claims, err := service.Authorize(&authorizationHeader)
	if err != nil {
		if err.Err == http.ErrNoCookie {
			return nil, &model.Erro{Err: service.NoAuthenticationToken, HttpCode: err.HttpCode}
		}
		return nil, err
	}

	if clientID == nil {
		return &claims.Id, nil
	}

	account, err := Services.AccountService.Get(clientID)
	if err != nil {
		return nil, err
	}

	if account.User_id != claims.Id {
		log.Error().Msg("User ID does not match with accounts User ID")
		return nil, &model.Erro{Err: errors.New("No match for user id"), HttpCode: http.StatusForbidden}
	}

	return nil, nil
}
