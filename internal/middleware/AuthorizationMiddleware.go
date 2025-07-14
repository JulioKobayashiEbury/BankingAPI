package middleware

import (
	"BankingAPI/internal/model"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

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
		if c.Path() == "/docs/*" || c.Path() == "/swagger.yaml" {
			next(c)
		}
		if c.Request().Header.Get(echo.HeaderAuthorization) == "" && c.FormValue("user_id") != "" {
			if c.Path() == "/auth/token" {
				return next(c)
			}
			return c.JSON(model.ErrNotAuthenticated.Code, model.ErrNotAuthenticated.Error())
		}
		authorizationHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authorizationHeader == "" {
			log.Warn().Msg("authorization header is missing")
			return c.JSON(model.ErrNotAuthenticated.Code, model.ErrNotAuthenticated.Error())
		}

		claims, err := service.Authorize(&authorizationHeader)
		if err != nil {
			log.Error().Msg("authorization failed: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		userResponse, err := h.userService.Get(c.Request().Context(), &claims.Id)
		if err != nil {
			log.Error().Msg("failed to get user: " + err.Error())
			return c.JSON(err.Code, err.Error())
		}

		if userResponse.Name == "admin" {
			log.Info().Msg("admin user authorized")
			return next(c)
		}

		if c.Param("user_id") != "" {
			if c.Param("user_id") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(model.ErrUserIDNotMatch.Code, model.ErrUserIDNotMatch.Internal)
			}
			return next(c)
		}
		if c.FormValue("user_id") != "" {
			if c.FormValue("user_id") != claims.Id {
				log.Error().Msg("user id doesn't match with claims id")
				return c.JSON(model.ErrUserIDNotMatch.Code, model.ErrUserIDNotMatch.Internal)
			}
			return next(c)
		}
		return next(c)
	}
}
