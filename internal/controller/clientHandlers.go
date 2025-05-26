package controller

import (
	"errors"
	"net/http"

	"BankingAPI/internal/service"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddClientsEndPoints(server *echo.Echo) {
	server.POST("/clients", ClientPostHandler)
	server.PUT("/clients/auth", ClientAuthHandler)
	server.GET("/clients/:client_id", ClientGetHandler)
	server.GET("/clients/:client_id/report", ClientGetReportHandler)
	server.DELETE("/clients/:client_id", ClientDeleteHandler)
	server.PUT("/clients/:client_id", ClientPutHandler)
	server.PUT("/clients/:client_id/block", ClientPutBlockHandler)
	server.PUT("/clients/:client_id/unblock", ClientPutUnBlockHandler)
}

func ClientPostHandler(c echo.Context) error {
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	clientResponse, err := service.CreateClient(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientAuthHandler(c echo.Context) error {
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ok, err := service.Authenticate(&(clientInfo).Client_id, &(clientInfo).Password, repository.ClientPath)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Credentials not valid")
	}
	cookie, err := service.GenerateToken(&(clientInfo.Client_id), service.ClientRole)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: "Client Authorized"})
}

func ClientGetHandler(c echo.Context) error {
	clientID, err := clientAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientInfo, err := service.Client(*clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientInfo))
}

func ClientDeleteHandler(c echo.Context) error {
	clientID, err := clientAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.ClientDelete(*clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client deleted seccesfully"})
}

func ClientPutHandler(c echo.Context) error {
	clientID, err := clientAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var clientInfo model.ClientRequest
	if err := c.Bind(&clientInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	clientInfo.Client_id = *clientID

	clientResponse, err := service.UpdateClient(&clientInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientResponse))
}

func ClientPutBlockHandler(c echo.Context) error {
	clientID, err := clientAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.ClientBlock(*clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Blocked"})
}

func ClientPutUnBlockHandler(c echo.Context) error {
	clientID, err := clientAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.ClientUnBlock(*clientID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Client Unblocked"})
}

func ClientGetReportHandler(c echo.Context) error {
	clientID, err := clientAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	clientReport, err := service.GenerateReportByClient(clientID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*clientReport))
}

func clientAuthorization(c *echo.Context) (*string, *model.Erro) {
	claims, err, cookie := service.Authorize((*c).Cookie("Token"))
	if err != nil {
		return nil, err
	}
	if cookie != nil {
		(*c).SetCookie(cookie)
	}
	clientID := (*c).Param("client_id")
	if (*claims).Id != clientID {
		return nil, &model.Erro{Err: errors.New("Not authorized"), HttpCode: http.StatusUnauthorized}
	}
	return &clientID, nil
}
