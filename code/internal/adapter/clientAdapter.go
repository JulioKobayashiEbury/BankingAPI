package adapter

import (
	"net/http"
	"strconv"

	adapter "BankingAPI/code/internal/adapter/infrastructure"
	"BankingAPI/code/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
)

func ClientPostAdapter(c *echo.Context) (int, interface{}) {
	var clientInfo domain.Client
	UserID, err := strconv.Atoi((*c).FormValue("UserID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	clientInfo.SetUserId(int32(UserID))
	clientInfo.SetName((*c).FormValue("Name"))
	clientInfo.SetDocument((*c).FormValue("Document"))
	clientInfo.SetPassword((*c).FormValue("Password"))
	clientInfo.SetStatus(true)

	valid := validator.New()
	if err := valid.Struct(clientInfo); err != nil {
		return http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid parameter"}
	}

	err = adapter.CreateClientDB(&clientInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, clientInfo
}

func ClientGetAdapter(c *echo.Context) (int, interface{}) {
	var clientInfo domain.Client

	clientID, err := strconv.Atoi((*c).Param("ClientID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	clientInfo.SetClientId(int32(clientID))

	err = adapter.GetClientDB(&clientInfo)

	return http.StatusOK, clientInfo
}

func ClientDeleteAdapter(c *echo.Context) (int, interface{}) {
	clientID, err := strconv.Atoi((*c).Param("ClientID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	deletedID, err := adapter.DeleteClientDB(int32(clientID))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, domain.Client{ClientId: deletedID}
}

func ClientPutAdapter(c *echo.Context) (int, interface{}) {
	var clientInfo domain.Client
	clientID, err := strconv.Atoi((*c).Param("ClientID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	clientInfo.SetClientId(int32(clientID))
	clientInfo.SetName((*c).FormValue("Name"))
	clientInfo.SetDocument((*c).FormValue("Document"))
	clientInfo.SetPassword((*c).FormValue("Password"))

	err = adapter.UpdateClientDB(&clientInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	return http.StatusOK, clientInfo
}

func ClientPutBlockAdapter(c *echo.Context) (int, interface{}) {
	var clientInfo domain.Client
	clientID, err := strconv.Atoi((*c).Param("ClientID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	clientInfo.SetClientId(int32(clientID))
	clientInfo.SetStatus(false)

	err = adapter.UpdateClientDB(&clientInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, clientInfo
}

func ClientPutUnBlockAdapter(c *echo.Context) (int, interface{}) {
	var clientInfo domain.Client
	clientID, err := strconv.Atoi((*c).Param("ClientID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	clientInfo.SetClientId(int32(clientID))
	clientInfo.SetStatus(true)

	err = adapter.UpdateClientDB(&clientInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, clientInfo
}
