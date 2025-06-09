package service

import (
	"errors"
	"net/http"
	"strings"

	model "BankingAPI/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var NoAuthenticationToken = errors.New("Not authenticated")

// validate token and authorize access to endpoint
func Authorize(authHeader *string) (*model.Claims, *model.Erro) {
	tokenString := strings.Split(*authHeader, " ")
	token := tokenString[len(tokenString)-1]
	if token == "" {
		log.Warn().Msg("No authentication token provided")
		return nil, &model.Erro{Err: NoAuthenticationToken, HttpCode: http.StatusUnauthorized}
	}
	var claims model.Claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, &model.Erro{Err: err, HttpCode: http.StatusUnauthorized}
		}
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	if !jwtToken.Valid {
		return nil, &model.Erro{Err: errors.New("Token not valid"), HttpCode: http.StatusUnauthorized}
	}

	log.Info().Msg("Authorized entrance: " + claims.Id)

	return &claims, nil
}
