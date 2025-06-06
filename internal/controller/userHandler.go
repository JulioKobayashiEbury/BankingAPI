package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/user"
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
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if (len(userInfo.Document) != documentLenghtIdeal) || (len(userInfo.Name) > maxNameLenght) {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	userResponse, err := Services.UserService.Create(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusCreated, (*userResponse))
}

func UserAuthHandler(c echo.Context) error {
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userDatabase := user.NewUserFireStore(DatabaseClient)
	serviceAuthenticate := service.NewAuth(userDatabase)

	ok, err := serviceAuthenticate.Authenticate(&(userInfo).User_id, &(userInfo).Password, model.UsersPath)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Credentials not valid")
	}

	cookie, err := serviceAuthenticate.GenerateToken(&(userInfo.User_id))
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	// usar isso ao invés de cookie?
	// c.Response().Header().Set(echo.HeaderAuthorization, token)

	c.SetCookie(cookie)
	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: "User Authorized"})
}

func UserPutBlockHandler(c echo.Context) error {
	userID, err := internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := Services.UserService.Status(userID, false); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Blocked"})
}

func UserPutUnBlockHandler(c echo.Context) error {
	userID, err := internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := Services.UserService.Status(userID, true); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Unblocked"})
}

func UserPutHandler(c echo.Context) error {
	userID, err := internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo.User_id = *userID
	userResponse, err := Services.UserService.Update(&userInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userResponse))
}

func UserGetHandler(c echo.Context) error {
	userID, err := internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	userResponse, err := Services.UserService.Get(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*userResponse))
}

func UserDeleteHandler(c echo.Context) error {
	userID, err := internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := Services.UserService.Delete(userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}

func UserGetReportHandler(c echo.Context) error {
	userID, err := internalUserAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	userReport, err := Services.UserService.Report(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userReport))
}

func internalUserAuthorization(c *echo.Context) (*string, *model.Erro) {
	claims, err := service.Authorize((*c).Cookie("Token"))
	if err != nil {
		if err.Err == http.ErrNoCookie {
			return nil, &model.Erro{Err: service.NoAuthenticationToken, HttpCode: err.HttpCode}
		}
		return nil, err
	}
	var userID string

	log.Debug().Msg("Entering users path...")
	userID = (*c).Param("user_id")
	if (*claims).Id != userID {
		return nil, &model.Erro{Err: service.NoAuthenticationToken, HttpCode: http.StatusBadRequest}
	}

	return &userID, nil
}

func externalUserAuthorization(c *echo.Context) *model.Erro {
	_, err := service.Authorize((*c).Cookie("Token"))
	if err != nil {
		if err.Err == http.ErrNoCookie {
			return &model.Erro{Err: service.NoAuthenticationToken, HttpCode: err.HttpCode}
		}
		return err
	}
	return nil
}
