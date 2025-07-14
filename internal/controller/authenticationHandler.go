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
	var userAuthInfo authInfo
	if err := c.Bind(&userAuthInfo); err != nil {
		obj, ok := err.(*echo.HTTPError)
		if !ok {
			return c.JSON(http.StatusInternalServerError, "could not resolve error in bind function")
		}
		return c.JSON(obj.Code, obj.Internal.Error())
	}
	if ok, err := h.authenticationService.Authenticate(c.Request().Context(), &userAuthInfo.User_Id, &userAuthInfo.Password); err != nil {
		return c.JSON(err.Code, err.Error())
	} else {
		if !ok {
			log.Error().Msg("authentication failed for user: " + userAuthInfo.User_Id)
			return c.JSON(http.StatusUnauthorized, model.StandartResponse{Message: "authentication failed"})
		}
		token, err := h.authenticationService.GenerateToken(c.Request().Context(), &userAuthInfo.User_Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.StandartResponse{Message: err.Error()})
		}
		c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+*token)
		return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Authorized"})
	}
}
