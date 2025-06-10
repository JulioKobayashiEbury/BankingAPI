package controller

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type AccountHandler interface {
	AccountPostHandler(c echo.Context) error
	AccountGetHandler(c echo.Context) error
	AccountDeleteHandler(c echo.Context) error
	AccountPutHandler(c echo.Context) error
	AccountPutBlockHandler(c echo.Context) error
	AccountPutUnBlockHandler(c echo.Context) error
	AccountGetReportHandler(c echo.Context) error
}

type accountHandlerImpl struct {
	accountService service.AccountService
}

func NewAccountHandler(accountServe service.AccountService) AccountHandler {
	return accountHandlerImpl{
		accountService: accountServe,
	}
}

func AddAccountEndPoints(server *echo.Echo, h AccountHandler) {
	// server.GET("/accounts", AccountGetOrderFilterHandler)

	server.GET("/accounts/:account_id", h.AccountGetHandler)
	server.GET("/accounts/report/:account_id", h.AccountGetReportHandler)
	server.POST("/accounts", h.AccountPostHandler)
	server.DELETE("/accounts/:account_id", h.AccountDeleteHandler)
	server.PUT("/accounts/:account_id", h.AccountPutHandler)
	server.PUT("/accounts/:account_id/block", h.AccountPutBlockHandler)
	server.PUT("/accounts/:account_id/unblock", h.AccountPutUnBlockHandler)
}

func (h accountHandlerImpl) AccountPostHandler(c echo.Context) error {
	userID, err := h.authorizationForAccountEndpoints(&c, nil)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	var accountInfo account.Account
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	log.Debug().Msg("Account userID: " + accountInfo.User_id + " UserID:" + *userID)
	if *userID != accountInfo.User_id {
		log.Warn().Msg("User ID does not match with accounts User ID")
		return c.JSON(http.StatusForbidden, model.StandartResponse{Message: "User ID does not match with accounts User ID"})
	}

	if accountInfo.Agency_id == 0 {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	accountResponse, err := h.accountService.Create(&accountInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusCreated, accountResponse)
}

func (h accountHandlerImpl) AccountGetHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if _, err := h.authorizationForAccountEndpoints(&c, &accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	accountResponse, err := h.accountService.Get(&accountID)
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
func (h accountHandlerImpl) AccountDeleteHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if _, err := h.authorizationForAccountEndpoints(&c, &accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.accountService.Delete(&accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Deleted"})
}

func (h accountHandlerImpl) AccountPutHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if _, err := h.authorizationForAccountEndpoints(&c, &accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	var accountInfo account.Account
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	accountInfo.Account_id = accountID

	accountResponse, err := h.accountService.Update(&accountInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, accountResponse)
}

func (h accountHandlerImpl) AccountPutBlockHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if _, err := h.authorizationForAccountEndpoints(&c, &accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.accountService.Status(&accountID, false); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Blocked Sucesfully!"})
}

func (h accountHandlerImpl) AccountPutUnBlockHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if _, err := h.authorizationForAccountEndpoints(&c, &accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.accountService.Status(&accountID, true); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Unblocked"})
}

func (h accountHandlerImpl) AccountGetReportHandler(c echo.Context) error {
	accountID := c.Param("account_id")

	if _, err := h.authorizationForAccountEndpoints(&c, &accountID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	accountReport, err := h.accountService.Report(&accountID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, accountReport)
}

func (h accountHandlerImpl) authorizationForAccountEndpoints(c *echo.Context, accountID *string) (*string, *model.Erro) {
	authorizationHeader := (*c).Request().Header.Get((echo.HeaderAuthorization))

	claims, err := service.Authorize(&authorizationHeader)
	if err != nil {
		return nil, err
	}

	if accountID == nil {
		return &claims.Id, nil
	}

	account, err := h.accountService.Get(accountID)
	if err != nil {
		return nil, err
	}
	if account.User_id != claims.Id {
		log.Error().Msg("User ID does not match with accounts User ID")
		return nil, &model.Erro{Err: errors.New("No match for user id"), HttpCode: http.StatusForbidden}
	}

	return nil, nil
}
