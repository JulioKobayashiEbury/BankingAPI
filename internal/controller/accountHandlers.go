package controller

import (
	"net/http"

	model "BankingAPI/internal/model/types"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddAccountEndPoints(server *echo.Echo) {
	server.GET("/accounts", AccountGetOrderFilterHandler)
	server.GET("/accounts/:account_id", AccountGetHandler)
	server.POST("/accounts", AccountPostHandler)
	server.DELETE("/accounts/:account_id", AccountDeleteHandler)
	server.PUT("/accounts/:account_id", AccountPutHandler)
	server.PUT("/accounts/:account_id/balance/withdrawal", AccountPutWithDrawalHandler)
	server.PUT("/accounts/:account_id/balance/deposit", AccountPutDepositHandler)
	server.PUT("/accounts/:account_id_from/newTransfer", AccountPostTranferHandler)
	server.PUT("/accounts/:account_id/block", AccountPutBlockHandler)
	server.PUT("/accounts/:account_id/unblock", AccountPutUnBlockHandler)
}

func AccountPostHandler(c echo.Context) error {
	var accountInfo model.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	accountResponse, err := service.CreateAccount(&accountInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, accountResponse)
}

func AccountGetHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	accountResponse, err := service.Account(accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, *accountResponse)
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	var listRequest model.ListRequest
	if err := c.Bind(&listRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	listOfAccounts, err := service.GetAccountByFilterAndOrder(&listRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, (*listOfAccounts))
}

func AccountDeleteHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if err := service.AccountDelete(accountID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Deleted"})
}

func AccountPutHandler(c echo.Context) error {
	var accountInfo model.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	accountInfo.Account_id = c.Param("account_id")

	accountResponse, err := service.UpdateAccount(&accountInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, accountResponse)
}

func AccountPutDepositHandler(c echo.Context) error {
	var depositRequest model.DepositRequest
	if err := c.Bind(&depositRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	depositRequest.Account_id = c.Param("account_id")

	newBalance, err := service.ProcessDeposit(&depositRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.DepositResponse{Account_id: depositRequest.Account_id, Balance: (*newBalance)})
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	var withdrawalRequest model.WithdrawalRequest
	if err := c.Bind(&withdrawalRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	withdrawalRequest.Account_id = c.Param("account_id")

	// talk to service
	newBalance, err := service.ProcessWithdrawal(&withdrawalRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.WithdrawalResponse{Account_id: withdrawalRequest.Account_id, Balance: (*newBalance)})
}

func AccountPutBlockHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if err := service.AccountBlock(accountID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Blocked Sucesfully!"})
}

func AccountPutUnBlockHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if err := service.AccountUnBlock(accountID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Unblocked"})
}

func AccountPutAutomaticDebit(c echo.Context) error {
	var newAutoDebit model.AutomaticDebitRequest
	if err := c.Bind(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	newAutoDebit.Account_id = c.Param("account_id")

	if err := service.ProcessNewAutomaticDebit(&newAutoDebit); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "New automatic debit is registered"})
}

func AccountPostTranferHandler(c echo.Context) error {
	var newTransferInfo model.TransferRequest
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	newTransferInfo.Account_id_from = c.Param("account_id_from")
	if err := service.ProcessNewTransfer(&newTransferInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Transfer was succesful"})
}
