package middleware

import (
	"errors"
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

var NotAuthenticated = &model.Erro{Err: errors.New("authorization header is missing, please authenticate first"), HttpCode: http.StatusUnauthorized}

type AuthMiddleware interface {
	AuthorizeMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddlewareImpl struct {
	userService           service.UserService
	authenticationService service.Authentication
}

func NewUserAuthMiddleware(userServe service.UserService) AuthMiddleware {
	return authMiddlewareImpl{
		userService: userServe,
	}
}

func (h authMiddlewareImpl) AuthorizeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderAuthorization) == "" && c.FormValue("user_id") != "" {
			if c.Path() == "/auth/token" {
				return next(c)
			}
			return c.JSON(NotAuthenticated.HttpCode, NotAuthenticated.Err.Error())
		}
		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authorizationHeader == "" {
			log.Warn().Msg("authorization header is missing")
			return c.JSON(NotAuthenticated.HttpCode, NotAuthenticated.Err.Error())
		}

		claims, err := service.Authorize(&authorizationHeader)
		if err != nil {
			log.Error().Msg("authorization failed: " + err.Err.Error())
			return c.JSON(err.HttpCode, err.Err.Error())
		}

		userResponse, err := h.userService.Get(&claims.Id)
		if err != nil {
			log.Error().Msg("failed to get user: " + err.Err.Error())
			return c.JSON(err.HttpCode, err.Err.Error())
		}

		if userResponse.Name == "admin" {
			log.Info().Msg("admin user authorized")
			return next(c)
		}

		if c.Param("user_id") != "" {
			if c.Param("user_id") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(http.StatusUnauthorized, model.StandartResponse{Message: "user id does not match with claims id"})
			}
			return next(c)
		}
		if c.FormValue("user_id") != "" {
			if c.FormValue("user_id") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(http.StatusUnauthorized, model.StandartResponse{Message: "user id does not match with claims id"})
			}
			return next(c)
		}
		return next(c)
	}
}
