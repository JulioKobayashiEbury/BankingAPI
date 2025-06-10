package controller

import (
	"errors"
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
	UserPutBlockHandler(c echo.Context) error
	UserPutUnBlockHandler(c echo.Context) error
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
	server.PUT("/users/:user_id/block", h.UserPutBlockHandler)
	server.PUT("/users/:user_id/unblock", h.UserPutUnBlockHandler)
}

func (h userHanderImpl) UserPostHandler(c echo.Context) error {
	if _, err := h.internalUserAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if (len(userInfo.Document) != documentLenghtIdeal) || (len(userInfo.Name) > maxNameLenght) {
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
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ok, err := h.authenticateService.Authenticate(&(userInfo).User_id, &(userInfo).Password, model.UsersPath)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Credentials not valid")
	}

	tokenString, err := h.authenticateService.GenerateToken(&(userInfo.User_id))
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	// usar isso ao invés de cookie?
	// c.Response().Header().Set(echo.HeaderAuthorization, token)

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+*tokenString)
	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: "User Authorized"})
}

func (h userHanderImpl) UserPutBlockHandler(c echo.Context) error {
	userID, err := h.internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.userService.Status(userID, false); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Blocked"})
}

func (h userHanderImpl) UserPutUnBlockHandler(c echo.Context) error {
	userID, err := h.internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.userService.Status(userID, true); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Unblocked"})
}

func (h userHanderImpl) UserPutHandler(c echo.Context) error {
	userID, err := h.internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo.User_id = *userID
	userResponse, err := h.userService.Update(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userResponse))
}

func (h userHanderImpl) UserGetHandler(c echo.Context) error {
	userID, err := h.internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	userResponse, err := h.userService.Get(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*userResponse))
}

func (h userHanderImpl) UserDeleteHandler(c echo.Context) error {
	userID, err := h.internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.userService.Delete(userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}

func (h userHanderImpl) UserGetReportHandler(c echo.Context) error {
	userID, err := h.internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	userReport, err := h.userService.Report(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userReport))
}

func (h userHanderImpl) internalUserAuthorization(c *echo.Context) (*string, *model.Erro) {
	authorizationHeader := (*c).Request().Header.Get(echo.HeaderAuthorization)

	claims, err := service.Authorize(&authorizationHeader)
	if err != nil {
		return nil, err
	}
	var userID string

	log.Debug().Msg("Entering users path...")

	userResponse, err := h.userService.Get(&claims.Id)
	if err != nil {
		return nil, err
	}
	userID = (*c).Param("user_id")
	if userResponse.Name == "admin" {
		return &userID, nil
	}

	if (*claims).Id != userID {
		log.Error().Msg("User ID does not match with accounts User ID")
		return nil, &model.Erro{Err: errors.New("No match for user id"), HttpCode: http.StatusForbidden}
	}

	return &userID, nil
}
