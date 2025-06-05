package controller

import (
	"net/http"

	model "BankingAPI/internal/model"
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

func AddAccountEndPoints(server *echo.Echo) {
	// server.GET("/accounts", AccountGetOrderFilterHandler)
	server.GET("/accounts/:account_id", AccountGetHandler)
	server.GET("/accounts/report/:account_id", AccountGetReportHandler)
	server.POST("/accounts", AccountPostHandler)
	server.DELETE("/accounts/:account_id", AccountDeleteHandler)
	server.PUT("/accounts/:account_id", AccountPutHandler)
	server.PUT("/accounts/:account_id/balance/withdrawal", AccountPutWithDrawalHandler)
	server.PUT("/accounts/:account_id/balance/deposit", AccountPutDepositHandler)
	server.PUT("/accounts/:account_id/transfer", AccountPutTransferHandler)
	server.PUT("/accounts/:account_id/block", AccountPutBlockHandler)
	server.PUT("/accounts/:account_id/unblock", AccountPutUnBlockHandler)
	server.PUT("/accounts/:account_id/debit", AccountPutAutomaticDebit)
}

func AccountPostHandler(c echo.Context) error {
	var accountInfo account.Account
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if accountInfo.Agency_id == 0 {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}
	userDatabase := user.NewUserFireStore(DatabaseClient)
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	accountDatabase := account.NewAccountFirestore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, clientDatabase, userDatabase)
	serviceCreate := service.NewCreateService(accountDatabase, clientDatabase, userDatabase, serviceGet)

	accountResponse, err := serviceCreate.CreateAccount(&accountInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusCreated, accountResponse)
}

func AccountGetHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	accountID := c.Param("account_id")

	userDatabase := user.NewUserFireStore(DatabaseClient)
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	accountDatabase := account.NewAccountFirestore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, clientDatabase, userDatabase)
	accountResponse, err := serviceGet.Account(accountID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *accountResponse)
}

/*
	func AccountGetOrderFilterHandler(c echo.Context) error {
		if _, err := userAuthorization(&c); err != nil {
		}
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
*/
func AccountDeleteHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	accountID := c.Param("account_id")
	serviceDelete := service.NewDeleteService(nil, nil, account.NewAccountFirestore(DatabaseClient))
	if err := serviceDelete.AccountDelete(accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Deleted"})
}

func AccountPutHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var accountInfo account.Account
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	accountInfo.Account_id = c.Param("account_id")

	userDatabase := user.NewUserFireStore(DatabaseClient)
	clientDatabase := client.NewClientFirestore(DatabaseClient)
	accountDatabase := account.NewAccountFirestore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, clientDatabase, userDatabase)
	serviceUpdate := service.NewUpdateService(accountDatabase, clientDatabase, userDatabase, serviceGet)
	accountResponse, err := serviceUpdate.UpdateAccount(&accountInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, accountResponse)
}

func AccountPutDepositHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var depositRequest deposit.Deposit
	if err := c.Bind(&depositRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	depositRequest.Account_id = c.Param("account_id")

	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	depositDatabase := deposit.NewDepositFirestore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, nil, nil)
	serviceDeposit := service.NewDepositService(accountDatabase, depositDatabase, serviceGet)

	depositID, err := serviceDeposit.ProcessDeposit(&depositRequest)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusAccepted, model.StandartResponse{Message: *depositID})
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var withdrawalRequest withdrawal.Withdrawal
	if err := c.Bind(&withdrawalRequest); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	withdrawalRequest.Account_id = c.Param("account_id")
	// talk to service
	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	withdrawalDatabase := withdrawal.NewWithdrawalFirestore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, nil, nil)
	serviceWithdrawal := service.NewWithdrawalService(accountDatabase, withdrawalDatabase, serviceGet)
	withdrawalID, err := serviceWithdrawal.ProcessWithdrawal(&withdrawalRequest)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: *withdrawalID})
}

func AccountPutBlockHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	accountID := c.Param("account_id")

	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	serviceGet := service.NewGetService(accountDatabase, nil, nil)
	serviceBlock := service.NewStatusService(nil, nil, accountDatabase, serviceGet)
	if err := serviceBlock.AccountBlock(accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Blocked Sucesfully!"})
}

func AccountPutUnBlockHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	accountID := c.Param("account_id")

	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	serviceGet := service.NewGetService(accountDatabase, nil, nil)
	serviceUnblock := service.NewStatusService(nil, nil, accountDatabase, serviceGet)
	if err := serviceUnblock.AccountUnBlock(accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Unblocked"})
}

func AccountPutAutomaticDebit(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var newAutoDebit automaticdebit.AutomaticDebit
	if err := c.Bind(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	newAutoDebit.Account_id = c.Param("account_id")

	autodebitDatabase := automaticdebit.NewAutoDebitFirestore(DatabaseClient)
	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	withdrawalDatabase := withdrawal.NewWithdrawalFirestore(DatabaseClient)

	serviceGet := service.NewGetService(accountDatabase, nil, nil)
	serviceWithdrawal := service.NewWithdrawalService(accountDatabase, withdrawalDatabase, serviceGet)
	serviceAutodebit := service.NewAutoDebitImpl(autodebitDatabase, serviceWithdrawal)
	autodebitResponse, err := serviceAutodebit.ProcessNewAutomaticDebit(&newAutoDebit)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusAccepted, *autodebitResponse)
}

func AccountPutTransferHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	var newTransferInfo transfer.Transfer
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	newTransferInfo.Account_id = c.Param("account_id")

	accountDatabase := account.NewAccountFirestore(DatabaseClient)
	transferDatabase := transfer.NewTransferFirestore(DatabaseClient)
	serviceTransfer := service.NewTransferService(accountDatabase, transferDatabase)
	transferResponse, err := serviceTransfer.ProcessNewTransfer(&newTransferInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *transferResponse)
}

func AccountGetReportHandler(c echo.Context) error {
	if _, err := userAuthorization(&c); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	accountID := c.Param("account_id")

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

	report, err := serviceReport.GenerateReportByAccount(&accountID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, report)
}
