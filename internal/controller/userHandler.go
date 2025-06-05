package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/model/withdrawal"
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

	userDatabase := user.NewUserFireStore(DatabaseClient)

	//passar esse nil é correto? eu deveria instanciar outros FirestoreEntities?
	//se sim, posso instanciar globalmente no contexto do User? ou controller?
	serviceGet := service.NewGetService(nil, nil, userDatabase)
	serviceCreate := service.NewCreateService(nil, nil, userDatabase, serviceGet)

	userResponse, err := serviceCreate.CreateUser(&userInfo)
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
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	userDatabase := user.NewUserFireStore(DatabaseClient)

	serviceGet := service.NewGetService(nil, nil, userDatabase)
	serviceStatus := service.NewStatusService(userDatabase, nil, nil, serviceGet)

	if err := serviceStatus.UserBlock(*userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Blocked"})
}

func UserPutUnBlockHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	userDatabase := user.NewUserFireStore(DatabaseClient)

	serviceGet := service.NewGetService(nil, nil, userDatabase)
	serviceStatus := service.NewStatusService(userDatabase, nil, nil, serviceGet)

	if err := serviceStatus.UserUnBlock(*userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User Unblocked"})
}

func UserPutHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var userInfo user.User
	if err := c.Bind(&userInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userDatabase := user.NewUserFireStore(DatabaseClient)

	serviceGet := service.NewGetService(nil, nil, userDatabase)
	serviceUpdate := service.NewUpdateService(nil, nil, userDatabase, serviceGet)

	userInfo.User_id = *userID
	userResponse, err := serviceUpdate.UpdateUser(&userInfo)
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
	userDatabase := user.NewUserFireStore(DatabaseClient)
	serviceGet := service.NewGetService(nil, nil, userDatabase)

	userResponse, err := serviceGet.User(*userID)
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

	userDatabase := user.NewUserFireStore(DatabaseClient)
	serviceDelete := service.NewDeleteService(userDatabase, nil, nil)

	if err := serviceDelete.UserDelete(*userID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "User deleted"})
}

func UserGetReportHandler(c echo.Context) error {
	userID, err := userAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	autodebitDatabase := automaticdebit.NewAutoDebitFirestore(DatabaseClient)
	withdrawalDatabase := withdrawal.NewWithdrawalFirestore(DatabaseClient)
	depositDatabase := deposit.NewDepositFirestore(DatabaseClient)
	transferDatabase := transfer.NewTransferFirestore(DatabaseClient)
	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	userDatabase := user.NewUserFireStore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, clientDatabase, userDatabase)
	serviceGetAll := service.NewGetAllService(autodebitDatabase, withdrawalDatabase, depositDatabase, transferDatabase, accountDatabase, clientDatabase)

	serviceReport := service.NewReportService(serviceGet, serviceGetAll)

	userReport, err := serviceReport.GenerateReportByUser(userID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, (*userReport))
}

func userAuthorization(c *echo.Context) (*string, *model.Erro) {
	claims, err := service.Authorize((*c).Cookie("Token"))
	if err != nil {
		if err.Err == http.ErrNoCookie {
			return nil, &model.Erro{Err: service.NoAuthenticationToken, HttpCode: err.HttpCode}
		}
		return nil, err
	}
	var userID string
	if (*c).Path() == "/users/:user_id" {
		log.Debug().Msg("Entering users path...")
		userID = (*c).Param("user_id")
		if (*claims).Id != userID {
			return nil, &model.Erro{Err: service.NoAuthenticationToken, HttpCode: http.StatusBadRequest}
		}
	}

	return &userID, nil
}
