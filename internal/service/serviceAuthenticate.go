package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

const expirationMin = 30

var jwtKey = []byte("bankingapi-key")

type auth struct {
	userDatabase user.UserRepository
}

func NewAuth(userDB user.UserRepository) Authentication {
	return auth{
		userDatabase: userDB,
	}
}

func (a auth) Authenticate(ctx context.Context, typeID *string, password *string) (bool, *echo.HTTPError) {
	userAuth, err := a.userDatabase.Get(ctx, typeID)
	if err != nil {
		return false, err
	}
	if *password != userAuth.Password {
		return false, &echo.HTTPError{Internal: errors.New("password is wrong"), Code: http.StatusBadRequest, Message: "password is wrong"}
	}
	return true, nil
}

func (a auth) GenerateToken(ctx context.Context, typeID *string) (*string, *echo.HTTPError) {
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
		return nil, &echo.HTTPError{Internal: err, Code: http.StatusInternalServerError, Message: err.Error()}
	}
	log.Info().Msg("Authenticated entrance: " + *typeID)

	return &tokenString, nil
}
