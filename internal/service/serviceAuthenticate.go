package service

import (
	"errors"
	"net/http"
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const expirationMin = 30

var jwtKey = []byte("bankingapi-key")

type Authentication interface {
	Authenticate(typeID *string, password *string, collection string) (bool, *model.Erro)
	GenerateToken(typeID *string) (*string, *model.Erro)
}

type auth struct {
	userDatabase model.RepositoryInterface
}

func NewAuth(userDB model.RepositoryInterface) Authentication {
	return auth{
		userDatabase: userDB,
	}
}

func (a auth) Authenticate(typeID *string, password *string, collection string) (bool, *model.Erro) {
	obj, err := a.userDatabase.Get(typeID)
	if err != nil {
		return false, err
	}

	userAuth, ok := obj.(*user.User)
	if !ok {
		return false, model.DataTypeWrong
	}

	if *password != userAuth.Password {
		return false, &model.Erro{Err: errors.New("Password is wrong"), HttpCode: http.StatusBadRequest}
	}
	return true, nil
}

func (a auth) GenerateToken(typeID *string) (*string, *model.Erro) {
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
	log.Info().Msg("Authenticated entrance: " + *typeID)

	return &tokenString, nil
}
