package controller

import (
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthenticationHandler interface {
	PostAuthenticationHandler(c echo.Context) error
}

type authInfo struct {
	User_Id  string `json:"user_id" xml:"user_id"`
	Password string `json:"password" xml:"password"`
}

type authenticationHandlerImpl struct {
	authenticationService service.Authentication
}

func NewAuthenticationHandler(authServe service.Authentication) AuthenticationHandler {
	return authenticationHandlerImpl{
		authenticationService: authServe,
	}
}

func AddAuthenticationEndpoints(server *echo.Echo, authHandler AuthenticationHandler) {
	server.POST("/auth/token", authHandler.PostAuthenticationHandler)
}

func (h authenticationHandlerImpl) PostAuthenticationHandler(c echo.Context) error {
	id := c.FormValue("user_id")
	password := c.FormValue("password")
	if id == "" || password == "" {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "user ID and password are required"})
	}
	if ok, err := h.authenticationService.Authenticate(c.Request().Context(), &id, &password); err != nil {
		return c.JSON(err.Code, err.Message)
	} else {
		if !ok {
			log.Error().Msg("authentication failed for user: " + id)
			return c.JSON(http.StatusUnauthorized, model.StandartResponse{Message: "authentication failed"})
		}
		token, err := h.authenticationService.GenerateToken(c.Request().Context(), &id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.StandartResponse{Message: err.Error()})
		}
		c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+*token)
		return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Authorized"})
	}
}
