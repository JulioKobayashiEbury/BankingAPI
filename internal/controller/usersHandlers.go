package controller

import (
	"net/http"
	"strconv"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddUsersEndPoints(server *echo.Echo) {
	server.POST("/users", UserPostHandler)
	server.GET("/users/:users_id", UserGetHandler)
	server.DELETE("/users/:users_id", UserDeleteHandler)
	server.PUT("/users/:users_id", UserPutHandler)
	server.PUT("/users/:users_id/block", UserPutBlockHandler)
	server.PUT("/users/:users_id/unblock", UserPutUnBlockHandler)
}

func UserPostHandler(c echo.Context) error {
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	userResponse, err := service.CreateUser(&userInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userResponse)
}

func UserPutBlockHandler(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := service.UserBlock(uint32(userID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Blocked"})
}

func UserPutUnBlockHandler(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := service.UserUnBlock(uint32(userID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Unblocked"})
}

func UserPutHandler(c echo.Context) error {
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	userResponse, err := service.UpdateUser(&userInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userResponse)
}

func UserGetHandler(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	userResponse, err := service.GetUser(uint32(userID))
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, userResponse)
}

func UserDeleteHandler(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := service.UserDelete(uint32(userID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}
