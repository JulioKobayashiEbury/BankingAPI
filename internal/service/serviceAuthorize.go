package service

import (
	"errors"
	"net/http"
	"strings"

	model "BankingAPI/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// validate token and authorize access to endpoint
func Authorize(authHeader *string) (*model.Claims, *echo.HTTPError) {
	tokenString := strings.Split(*authHeader, " ")
	token := tokenString[len(tokenString)-1]
	if token == "" {
		log.Warn().Msg("No authentication token provided")
		return nil, model.ErrNotAuthenticated
	}
	var claims model.Claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusUnauthorized, Message: err.Error()}
	}
	if !jwtToken.Valid {
		return nil, &echo.HTTPError{Internal: errors.New("token not valid"), Code: http.StatusUnauthorized, Message: "token not valid"}
	}

	log.Info().Msg("Authorized entrance: " + claims.Id)

	return &claims, nil
}
