package controller

import (
	"errors"
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddUsersEndPoints(server *echo.Echo) {
	server.POST("/users", UserPostHandler)
	server.PUT("/users/auth", UserAuthHandler)
	server.GET("/users/:user_id", UserGetHandler)
	server.GET("/users/:user_id/report", UserGetReportHandler)
	server.DELETE("/users/:user_id", UserDeleteHandler)
	server.PUT("/users/:user_id", UserPutHandler)
	server.PUT("/users/:user_id/block", UserPutBlockHandler)
	server.PUT("/users/:user_id/unblock", UserPutUnBlockHandler)
}

func UserPostHandler(c echo.Context) error {
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	userResponse, err := service.CreateUser(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusCreated, (*userResponse))
}

func UserAuthHandler(c echo.Context) error {
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	ok, err := service.Authenticate(&(userInfo).User_id, &(userInfo).Password, repository.UsersPath)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Credentials not valid")
	}
	cookie, err := service.GenerateToken(&(userInfo.User_id), service.UserRole)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	// usar isso ao invés de cookie?
	// c.Response().Header().Set(echo.HeaderAuthorization, token)

	c.SetCookie(cookie)
	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: "User Authorized"})
}

func UserPutBlockHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.UserBlock(*userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Blocked"})
}

func UserPutUnBlockHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := service.UserUnBlock(*userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Unblocked"})
}

func UserPutHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	userInfo.User_id = *userID
	userResponse, err := service.UpdateUser(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userResponse))
}

func UserGetHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	userResponse, err := service.User(*userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*userResponse))
}

func UserDeleteHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.UserDelete(*userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}

func UserGetReportHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	userReport, err := service.GenerateReportByUser(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userReport))
}

func userAuthorization(c *echo.Context) (*string, *model.Erro) {
	claims, err, cookie := service.Authorize((*c).Cookie("Token"))
	if err != nil {
		return nil, err
	}
	if cookie != nil {
		(*c).SetCookie(cookie)
	}
	userID := (*c).Param("user_id")
	if (*claims).Id != userID {
		return nil, &model.Erro{Err: errors.New("Not authorized"), HttpCode: http.StatusBadRequest}
	}
	return &userID, nil
}
