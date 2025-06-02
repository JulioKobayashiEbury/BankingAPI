package service

import (
	"errors"
	"net/http"
	"time"

	model "BankingAPI/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

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

	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) < 30*time.Second {
		cookie, err := GenerateToken(&(claims.Id), claims.Role)
		return &claims, err, cookie
	}
	return &claims, nil, nil
}
