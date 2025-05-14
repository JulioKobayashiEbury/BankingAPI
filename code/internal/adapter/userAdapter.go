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

func UserPostAdapter(c echo.Context) (int, interface{}) {
	var userInfo domain.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg("Error binding struct")
		return http.StatusBadRequest, "Binding failed, review query"
	}
	valid := validator.New()
	if err := valid.Struct(userInfo); err != nil {
		return http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid parameter"}
	}
	return http.StatusCreated, application.CreateUser(&userInfo)
}

func UserGetAdapter(c echo.Context) (int, interface{}) {
	allUsers, err := application.GetAllUsers()
	if err != nil {
		log.Info().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: "Error getting all users"}
	}
	return http.StatusOK, allUsers
}

func UserDeleteAdapter(c echo.Context) (int, interface{}) {
	userID, err := strconv.Atoi(c.ParamValues()[0])
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

func UserPutAdapter(c echo.Context) (int, interface{}) {
	var userInfo *domain.User
	userID, err := strconv.Atoi(c.Param("UserID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, err.Error()
	}
	userInfo.SetId(int32(userID))
	userInfo.SetName(c.FormValue("Name"))
	userInfo.SetPassword(c.FormValue("Password"))
	userInfo.SetDocument(c.FormValue("Document"))

	userInfo, err = adapter.UpdateUserDB(userInfo)

	return http.StatusAccepted, userInfo
}
