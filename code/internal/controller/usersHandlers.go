package controller

import (
	"net/http"

	"BankingAPI/code/internal/domain"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
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
	var userInfo domain.UserRequest
	if err := c.Bind(&userInfo); err != nil {
	}
	valid := validator.New()
	if err := valid.Struct(userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid parameter"})
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserPutBlockHandler(c echo.Context) error {
	var userInfo domain.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserPutUnBlockHandler(c echo.Context) error {
	var userInfo domain.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserPutHandler(c echo.Context) error {
	var userInfo domain.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	return c.JSON(http.StatusOK, userInfo)
}

func UserGetHandler(c echo.Context) error {
	var userInfo domain.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, userInfo)
}

func UserDeleteHandler(c echo.Context) error {
	var userInfo domain.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, userInfo)
}
