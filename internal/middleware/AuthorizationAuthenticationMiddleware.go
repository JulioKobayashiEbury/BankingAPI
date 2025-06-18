package middleware

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/service"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware interface {
	AuthorizeMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddlewareImpl struct {
	userService           service.UserService
	authenticationService service.Authentication
}

func NewUserAuthMiddleware(userServe service.UserService, authServe service.Authentication) AuthMiddleware {
	return authMiddlewareImpl{
		userService:           userServe,
		authenticationService: authServe,
	}
}

func (h authMiddlewareImpl) AuthorizeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderAuthorization) == "" && c.FormValue("user_id") != "" && c.Path() == "/users/auth" {
			userID := c.FormValue("user_id")
			password := c.FormValue("password")
			if ok, err := h.authenticationService.Authenticate(&userID, &password); err != nil {
				return c.JSON(err.HttpCode, err.Err.Error())
			} else {
				if !ok {
					log.Error().Msg("authentication failed for user: " + userID)
					return c.JSON(http.StatusUnauthorized, model.StandartResponse{Message: "authentication failed"})
				}
				token, err := h.authenticationService.GenerateToken(&userID)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, model.StandartResponse{Message: err.Err.Error()})
				}
				c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+*token)
				return next(c)
			}
		}
		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authorizationHeader == "" {
			log.Warn().Msg("authorization header is missing")
			return c.JSON(http.StatusUnauthorized, model.StandartResponse{Message: "authorization header is missing"})
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
