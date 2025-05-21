package controller

import (
	"net/http"

	model "BankingAPI/internal/model/types"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddUsersEndPoints(server *echo.Echo) {
	server.POST("/users", UserPostHandler)
	server.POST("/users/auth", UserAuthHandler)
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
		return c.JSON(err.HttpCode, err.Err)
	}
	return c.JSON(http.StatusCreated, (*userResponse))
}

func UserAuthHandler(c echo.Context) error {
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	ok, err := service.AuthenticateUser(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if ok != true {
		return c.JSON(http.StatusUnauthorized, "Credentials not valid")
	}
	cookie, err := service.GenerateToken(&(userInfo.User_id), service.UserRole)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Authorized"})
}

func UserPutBlockHandler(c echo.Context) error {
	claims, err, cookie := service.Authorize(c.Cookie("Token"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Err.Error())
	}
	if cookie != nil {
		c.SetCookie(cookie)
	}
	userID := c.Param("users_id")
	if (*claims).Id != userID {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Not authorized"})
	}
	if err := service.UserBlock(userID); err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Blocked"})
}

func UserPutUnBlockHandler(c echo.Context) error {
	claims, err, cookie := service.Authorize(c.Cookie("Token"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Err.Error())
	}
	if cookie != nil {
		c.SetCookie(cookie)
	}
	userID := c.Param("users_id")
	if (*claims).Id != userID {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Not authorized"})
	}

	if err := service.UserUnBlock(userID); err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Unblocked"})
}

func UserPutHandler(c echo.Context) error {
	claims, err, cookie := service.Authorize(c.Cookie("Token"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Err.Error())
	}
	if cookie != nil {
		c.SetCookie(cookie)
	}
	userID := c.Param("users_id")
	if (*claims).Id != userID {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Not authorized"})
	}
	var userInfo model.UserRequest
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	userResponse, err := service.UpdateUser(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userResponse))
}

func UserGetHandler(c echo.Context) error {
	claims, err, cookie := service.Authorize(c.Cookie("Token"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Err.Error())
	}
	if cookie != nil {
		c.SetCookie(cookie)
	}
	userID := c.Param("users_id")
	if (*claims).Id != userID {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Not authorized"})
	}
	userResponse, err := service.User(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*userResponse))
}

func UserDeleteHandler(c echo.Context) error {
	claims, err, cookie := service.Authorize(c.Cookie("Token"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Err.Error())
	}
	if cookie != nil {
		c.SetCookie(cookie)
	}
	userID := c.Param("users_id")
	if (*claims).Id != userID {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Not authorized"})
	}
	if err := service.UserDelete(userID); err != nil {
		return c.JSON(err.HttpCode, err.Err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}
