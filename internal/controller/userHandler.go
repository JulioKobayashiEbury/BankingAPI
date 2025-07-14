package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	UserPostHandler(c echo.Context) error
	UserPutHandler(c echo.Context) error
	UserGetHandler(c echo.Context) error
	UserDeleteHandler(c echo.Context) error
	UserGetReportHandler(c echo.Context) error
}

type userHandlerImpl struct {
	userService service.UserService
}

func NewUserHandler(userServe service.UserService, authServe service.Authentication) UserHandler {
	return &userHandlerImpl{
		userService: userServe,
	}
}

func AddUsersEndPoints(group *echo.Group, h UserHandler) {
	group.POST("/users", h.UserPostHandler)
	group.GET("/users/:user_id", h.UserGetHandler)
	group.GET("/users/:user_id/report", h.UserGetReportHandler)
	group.DELETE("/users/:user_id", h.UserDeleteHandler)
	group.PUT("/users/:user_id", h.UserPutHandler)
}

func (h userHandlerImpl) UserPostHandler(c echo.Context) error {
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if (len(userInfo.Document) != documentLenghtForUser) || (len(userInfo.Name) > maxNameLenght) {
		log.Warn().Msg("User parameters are not ideal for creating user")
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	userResponse, err := h.userService.Create(c.Request().Context(), &userInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusCreated, (*userResponse))
}

func (h userHandlerImpl) UserPutHandler(c echo.Context) error {
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userID := c.Param("user_id")

	userInfo.User_id = userID

	userResponse, err := h.userService.Update(c.Request().Context(), &userInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusOK, (*userResponse))
}

func (h userHandlerImpl) UserGetHandler(c echo.Context) error {
	userID := c.Param("user_id")

	userResponse, err := h.userService.Get(c.Request().Context(), &userID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, (*userResponse))
}

func (h userHandlerImpl) UserDeleteHandler(c echo.Context) error {
	userID := c.Param("user_id")

	if err := h.userService.Delete(c.Request().Context(), &userID); err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}

func (h userHandlerImpl) UserGetReportHandler(c echo.Context) error {
	userID := c.Param("user_id")

	userReport, err := h.userService.Report(c.Request().Context(), &userID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusOK, (*userReport))
}
