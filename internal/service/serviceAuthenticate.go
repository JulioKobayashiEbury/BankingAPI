package service

import (
	"net/http"
	"time"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const expirationMin = 30

var jwtKey = []byte("bankingapi-key")

func Authenticate(typeID *string, password *string, collection string) (bool, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(typeID, collection)
	if err != nil {
		return false, err
	}
	var typeAuth Auth
	if err := docSnapshot.DataTo(&typeAuth); err != nil {
		return false, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	if *password != typeAuth.Password {
		return false, nil
	}
	log.Info().Msg("Authenticated entrance: " + *typeID)
	return true, nil
}

func GenerateToken(typeID *string, role string) (*http.Cookie, *model.Erro) {
	expirationTime := time.Now().Add(time.Minute * expirationMin)
	Claim := &model.Claims{
		Id:   (*typeID),
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *Claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	return &http.Cookie{
		Name:    "Token",
		Value:   tokenString,
		Expires: expirationTime,
	}, nil
}
