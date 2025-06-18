package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	UserPostHandler(c echo.Context) error
	UserAuthHandler(c echo.Context) error
	UserPutHandler(c echo.Context) error
	UserGetHandler(c echo.Context) error
	UserDeleteHandler(c echo.Context) error
	UserGetReportHandler(c echo.Context) error
}

type userHanderImpl struct {
	userService         service.UserService
	authenticateService service.Authentication
}

func NewUserHandler(userServe service.UserService, authServe service.Authentication) UserHandler {
	return &userHanderImpl{
		userService:         userServe,
		authenticateService: authServe,
	}
}

func AddUsersEndPoints(server *echo.Echo, h UserHandler) {
	server.POST("/users", h.UserPostHandler)
	server.PUT("/users/auth", h.UserAuthHandler)
	server.GET("/users/:user_id", h.UserGetHandler)
	server.GET("/users/:user_id/report", h.UserGetReportHandler)
	server.DELETE("/users/:user_id", h.UserDeleteHandler)
	server.PUT("/users/:user_id", h.UserPutHandler)
}

func (h userHanderImpl) UserPostHandler(c echo.Context) error {
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if (len(userInfo.Document) != documentLenghtForUser) || (len(userInfo.Name) > maxNameLenght) {
		log.Warn().Msg("User parameters are not ideal for creating user")
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	userResponse, err := h.userService.Create(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusCreated, (*userResponse))
}

func (h userHanderImpl) UserAuthHandler(c echo.Context) error {
	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: "User Authorized"})
}

func (h userHanderImpl) UserPutHandler(c echo.Context) error {
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userID := c.Param("user_id")

	userInfo.User_id = userID

	userResponse, err := h.userService.Update(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userResponse))
}

func (h userHanderImpl) UserGetHandler(c echo.Context) error {
	userID := c.Param("user_id")

	userResponse, err := h.userService.Get(&userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*userResponse))
}

func (h userHanderImpl) UserDeleteHandler(c echo.Context) error {
	userID := c.Param("user_id")

	if err := h.userService.Delete(&userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}

func (h userHanderImpl) UserGetReportHandler(c echo.Context) error {
	userID := c.Param("user_id")

	userReport, err := h.userService.Report(&userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userReport))
}
