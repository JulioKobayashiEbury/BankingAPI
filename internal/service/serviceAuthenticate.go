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

// validate user identity and generate token
func AuthenticateUser(userInfo *model.UserRequest) (bool, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&(userInfo.User_id), repository.UsersPath)
	if err != nil {
		return false, err
	}
	var userAuth Auth
	if err := docSnapshot.DataTo(&userAuth); err != nil {
		return false, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	if (*userInfo).Password != userAuth.Password {
		return false, nil
	}
	log.Info().Msg("Authenticated entrance: " + (*&userInfo.User_id))
	return true, nil
}

func AuthenticateClient(clientInfo *model.ClientRequest) (bool, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&(clientInfo.Client_id), repository.ClientPath)
	if err != nil {
		return false, err
	}
	var clientAuth Auth
	if err := docSnapshot.DataTo(&clientAuth); err != nil {
		return false, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	if (*clientInfo).Password != clientAuth.Password {
		return false, nil
	}
	log.Info().Msg("Authenticated entrance: " + (*&clientInfo.Client_id))
	return true, nil
}

func AuthenticateAccount(accountInfo *model.AccountRequest) (bool, *model.Erro) {
	docSnapshot, err := repository.GetTypeFromDB(&(accountInfo.Account_id), repository.AccountsPath)
	if err != nil {
		return false, err
	}
	var accountAuth Auth
	if err := docSnapshot.DataTo(&accountAuth); err != nil {
		return false, &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	if (*accountInfo).Password != accountAuth.Password {
		return false, nil
	}
	log.Info().Msg("Authenticated entrance: " + (*&accountInfo.Account_id))
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
