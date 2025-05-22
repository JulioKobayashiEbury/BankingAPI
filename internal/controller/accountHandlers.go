package controller

import (
	"errors"
	"net/http"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddAccountEndPoints(server *echo.Echo) {
	server.GET("/accounts", AccountGetOrderFilterHandler)
	server.GET("/accounts/:account_id", AccountGetHandler)
	server.POST("/accounts", AccountPostHandler)
	server.PUT("/accounts/auth", AccountAuthHandler)
	server.DELETE("/accounts/:account_id", AccountDeleteHandler)
	server.PUT("/accounts/:account_id", AccountPutHandler)
	server.PUT("/accounts/:account_id/balance/withdrawal", AccountPutWithDrawalHandler)
	server.PUT("/accounts/:account_id/balance/deposit", AccountPutDepositHandler)
	server.PUT("/accounts/:account_id/newTransfer", AccountPutTransferHandler)
	server.PUT("/accounts/:account_id/block", AccountPutBlockHandler)
	server.PUT("/accounts/:account_id/unblock", AccountPutUnBlockHandler)
	server.PUT("/accounts/:account_id/debit", AccountPutAutomaticDebit)
}

func AccountPostHandler(c echo.Context) error {
	var accountInfo model.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	accountResponse, err := service.CreateAccount(&accountInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusCreated, accountResponse)
}

func AccountAuthHandler(c echo.Context) error {
	var accountInfo model.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	ok, err := service.Authenticate(&(accountInfo).Account_id, &(accountInfo).Password, repository.AccountsPath)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Credentials not valid")
	}
	cookie, err := service.GenerateToken(&(accountInfo.Account_id), service.AccountRole)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: "Account Authorized"})
}

func AccountGetHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	accountResponse, err := service.Account(*accountID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *accountResponse)
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	var listRequest model.ListRequest
	if err := c.Bind(&listRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	listOfAccounts, err := service.GetAccountByFilterAndOrder(&listRequest)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, (*listOfAccounts))
}

func AccountDeleteHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.AccountDelete(*accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Deleted"})
}

func AccountPutHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var accountInfo model.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	accountInfo.Account_id = *accountID

	accountResponse, err := service.UpdateAccount(&accountInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, accountResponse)
}

func AccountPutDepositHandler(c echo.Context) error {
	var depositRequest model.DepositRequest
	if err := c.Bind(&depositRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	depositRequest.Account_id = c.Param("account_id")

	newBalance, err := service.ProcessDeposit(&depositRequest)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.DepositResponse{Account_id: depositRequest.Account_id, Balance: (*newBalance)})
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var withdrawalRequest model.WithdrawalRequest
	if err := c.Bind(&withdrawalRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	withdrawalRequest.Account_id = *accountID

	// talk to service
	newBalance, err := service.ProcessWithdrawal(&withdrawalRequest)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.WithdrawalResponse{Account_id: withdrawalRequest.Account_id, Balance: (*newBalance)})
}

func AccountPutBlockHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.AccountBlock(*accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Blocked Sucesfully!"})
}

func AccountPutUnBlockHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := service.AccountUnBlock(*accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Unblocked"})
}

func AccountPutAutomaticDebit(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var newAutoDebit model.AutomaticDebitRequest
	if err := c.Bind(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	newAutoDebit.Account_id = *accountID

	if err := service.ProcessNewAutomaticDebit(&newAutoDebit); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "New automatic debit is registered"})
}

func AccountPutTransferHandler(c echo.Context) error {
	accountID, err := accountAuthorization(&c)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var newTransferInfo model.TransferRequest
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	newTransferInfo.Account_id = *accountID
	if err := service.ProcessNewTransfer(&newTransferInfo); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Transfer was succesful: TransferID:" + newTransferInfo.Transfer_id})
}

func accountAuthorization(c *echo.Context) (*string, *model.Erro) {
	claims, err, cookie := service.Authorize((*c).Cookie("Token"))
	if err != nil {
		return nil, err
	}
	if cookie != nil {
		(*c).SetCookie(cookie)
	}
	userID := (*c).Param("account_id")
	if (*claims).Id != userID {
		return nil, &model.Erro{Err: errors.New("Not authorized"), HttpCode: http.StatusBadRequest}
	}
	return &userID, nil
}
