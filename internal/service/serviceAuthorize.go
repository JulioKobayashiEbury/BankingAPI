package service

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var NoAuthenticationToken = errors.New("Not authenticated")

// validate token and authorize access to endpoint
func Authorize(cookie *http.Cookie, err error) (*model.Claims, *model.Erro, *http.Cookie) {
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}, nil
	}
	tokenString := cookie.Value
	var claims model.Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, &model.Erro{Err: err, HttpCode: http.StatusUnauthorized}, nil
		}
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}, nil
	}
	if !token.Valid {
		return nil, &model.Erro{Err: errors.New("Token not valid"), HttpCode: http.StatusUnauthorized}, nil
	}

	log.Info().Msg("Authorized entrance: " + claims.Id)

	return &claims, nil, nil
}
