package adapter

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"

	adapter "BankingAPI/code/internal/adapter/infrastructure"
	"BankingAPI/code/internal/application"
	"BankingAPI/code/internal/domain"
)

func UserPostAdapter(c *echo.Context) (int, interface{}) {
	var userInfo domain.User
	userInfo.SetName((*c).FormValue("Name"))
	userInfo.SetDocument((*c).FormValue("Document"))
	userInfo.SetPassword((*c).FormValue("Password"))
	userInfo.SetStatus(true)

	valid := validator.New()
	if err := valid.Struct(userInfo); err != nil {
		return http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid parameter"}
	}
	err := adapter.CreateUserDB(&userInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	return http.StatusCreated, userInfo
}

func UserGetAdapter(c *echo.Context) (int, interface{}) {
	allUsers, err := application.GetAllUsers()
	if err != nil {
		log.Info().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: "Error getting all users"}
	}
	return http.StatusOK, allUsers
}

func UserDeleteAdapter(c *echo.Context) (int, interface{}) {
	userID, err := strconv.Atoi((*c).ParamValues()[0])
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	// func to access db and delete user
	deletedID, err := adapter.DeleteUserDB(int32(userID))
	if err != nil {
		log.Warn().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	return http.StatusOK, domain.User{Id: deletedID}
}

func UserPutAdapter(c *echo.Context) (int, interface{}) {
	var userInfo domain.User
	userID, err := strconv.Atoi((*c).Param("UserID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, err.Error()
	}
	userInfo.SetId(int32(userID))
	userInfo.SetName((*c).FormValue("Name"))
	userInfo.SetPassword((*c).FormValue("Password"))
	userInfo.SetDocument((*c).FormValue("Document"))

	err = adapter.UpdateUserDB(&userInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusAccepted, userInfo
}

func UserPutBlock(c *echo.Context) (int, interface{}) {
	userID, err := strconv.Atoi((*c).Param("UserID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	var userInfo domain.User
	userInfo.SetId(int32(userID))
	userInfo.SetStatus(false)

	err = adapter.UpdateUserDB(&userInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	return http.StatusOK, userInfo
}

func UserPutUnblock(c *echo.Context) (int, interface{}) {
	userID, err := strconv.Atoi((*c).Param("UserID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	var userInfo domain.User
	userInfo.SetId(int32(userID))
	userInfo.SetStatus(true)

	err = adapter.UpdateUserDB(&userInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	return http.StatusOK, userInfo
}
