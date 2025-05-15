package controller

import (
	controller "BankingAPI/code/internal/controller/objects"
	"BankingAPI/code/internal/service"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddAccountEndPoints(server *echo.Echo) {
	server.GET("/accounts", AccountGetOrderFilterHandler)
	server.GET("/accounts/:account_id", AccountGetHandler)
	server.POST("/accounts", AccountPostHandler)
	server.DELETE("/accounts/:account_id", AccountDeleteHandler)
	server.PUT("/accounts/:account_id/withdrawal", AccountPutWithDrawalHandler)
	server.PUT("/accounts/:account_id/deposit", AccountPutDepositHandler)
}

func AccountPostHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountGetHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountDeleteHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutDepositHandler(c echo.Context) error {
	var depositRequest controller.DepositRequest
	if err := c.Bind(&depositRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	newBalance, err := service.ProcessDeposit(&depositRequest)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.DepositResponse{Account_id: depositRequest.Account_id, Balance: (*newBalance)})
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	var withdrawalRequest controller.WithdrawalRequest
	if err := c.Bind(&withdrawalRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service
	newBalance, err := service.ProcessWithdrawal(&withdrawalRequest)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.WithdrawalResponse{Account_id: withdrawalRequest.Account_id, Balance: (*newBalance)})
}

func AccountPutBlockHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := service.AccountBlock(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	// talk to service

	return c.JSON(http.StatusOK, controller.BlockUnBlockResponse{Message: "Account Blocked Sucesfully!"})
}

func AccountPutUnBlockHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutAutomaticDebit(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPostTranferHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}
