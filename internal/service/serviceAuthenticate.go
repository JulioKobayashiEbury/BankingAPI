package service

import (
	"net/http"
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const expirationMin = 30

var jwtKey = []byte("bankingapi-key")

func Authenticate(typeID *string, password *string, collection string) (bool, *model.Erro) {
	database := &user.UserFirestore{}
	database.Request = &user.UserRequest{
		User_id: *typeID,
	}
	if err := database.GetAuthInfo(); err != nil {
		return false, err
	}
	if *password != database.AuthUser.Password {
		return false, nil
	}
	log.Info().Msg("Authenticated entrance: " + *typeID)
	return true, nil
}

func GenerateToken(typeID *string, role string) (*http.Cookie, *model.Erro) {
	expirationTime := time.Now().Add(time.Minute * expirationMin)
	Claim := &model.Claims{
		Id: (*typeID),
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
