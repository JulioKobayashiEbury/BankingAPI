package service

import (
	"errors"
	"net/http"
	"strings"

	model "BankingAPI/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var NoAuthenticationToken = &model.Erro{Err: errors.New("not authenticated"), HttpCode: http.StatusUnauthorized}

// validate token and authorize access to endpoint
func Authorize(authHeader *string) (*model.Claims, *model.Erro) {
	tokenString := strings.Split(*authHeader, " ")
	token := tokenString[len(tokenString)-1]
	if token == "" {
		log.Warn().Msg("No authentication token provided")
		return nil, NoAuthenticationToken
	}
	var claims model.Claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusUnauthorized}
	}
	if !jwtToken.Valid {
		return nil, &model.Erro{Err: errors.New("token not valid"), HttpCode: http.StatusUnauthorized}
	}

	log.Info().Msg("Authorized entrance: " + claims.Id)

	return &claims, nil
}
