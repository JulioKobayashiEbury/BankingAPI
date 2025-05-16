package controller

import (
	controller "BankingAPI/code/internal/controller/objects"
	service "BankingAPI/code/internal/service"
	"net/http"
	"strconv"

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
	server.PUT("/accounts/:account_id/newTransfer", AccountPostTranferHandler)
	server.PUT("/accounts/:account_id/block", AccountPutBlockHandler)
	server.PUT("/accounts/:account_id/unblock", AccountPutUnBlockHandler)
}

func AccountPostHandler(c echo.Context) error {
	var accountInfo controller.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	accountResponse, err := service.CreateAccount(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, accountResponse)
}

func AccountGetHandler(c echo.Context) error {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	accountResponse, err := service.GetAccount(uint32(accountID))
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, *accountResponse)
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	var listRequest controller.ListRequest
	if err := c.Bind(&listRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	listOfAccounts, err := service.GetAccountByFilterAndOrder(&listRequest)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, (*listOfAccounts))
}

func AccountDeleteHandler(c echo.Context) error {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.DeleteAccount(uint32(accountID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.StandartResponse{Message: "Account Deleted"})
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

	accountID, err := strconv.ParseUint(c.Param("account_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := service.AccountBlock(uint32(accountID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.StandartResponse{Message: "Account Blocked Sucesfully!"})
}

func AccountPutUnBlockHandler(c echo.Context) error {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.AccountUnBlock(uint32(accountID)); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.StandartResponse{Message: "Account Unblocked"})
}

func AccountPutAutomaticDebit(c echo.Context) error {
	var newAutoDebit controller.AutomaticDebitRequest
	if err := c.Bind(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.ProcessNewAutomaticDebit(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.StandartResponse{Message: "New automatic debit is registered"})
}

func AccountPostTranferHandler(c echo.Context) error {
	var newTransferInfo controller.TransferRequest
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := service.ProcessNewTransfer(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, controller.StandartResponse{Message: "Transfer was succesful"})

}
