package controller

import (
	controller "BankingAPI/code/internal/controller/objects"
	"net/http"

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
	var userInfo controller.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserPutBlockHandler(c echo.Context) error {
	var userInfo controller.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserPutUnBlockHandler(c echo.Context) error {
	var userInfo controller.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserPutHandler(c echo.Context) error {
	var userInfo controller.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserGetHandler(c echo.Context) error {
	var userInfo controller.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, userInfo)
}

func UserDeleteHandler(c echo.Context) error {
	var userInfo controller.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, userInfo)
}
