package controller

import (
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AccountHandler interface {
	AccountPostHandler(c echo.Context) error
	AccountGetHandler(c echo.Context) error
	AccountDeleteHandler(c echo.Context) error
	AccountPutHandler(c echo.Context) error
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
	server.GET("/accounts/:account_id", h.AccountGetHandler)
	server.GET("/accounts/:account_id/report", h.AccountGetReportHandler)
	server.POST("/accounts", h.AccountPostHandler)
	server.DELETE("/accounts/:account_id", h.AccountDeleteHandler)
	server.PUT("/accounts/:account_id", h.AccountPutHandler)
}

func (h accountHandlerImpl) AccountPostHandler(c echo.Context) error {
	var accountInfo account.Account
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if !validateInput(nil, &accountInfo) {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	accountResponse, err := h.accountService.Create(c.Request().Context(), &accountInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusCreated, accountResponse)
}

func (h accountHandlerImpl) AccountGetHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if !validateInput(&accountID, nil) {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	accountResponse, err := h.accountService.Get(c.Request().Context(), &accountID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, *accountResponse)
}

func (h accountHandlerImpl) AccountDeleteHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if !validateInput(&accountID, nil) {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	if err := h.accountService.Delete(c.Request().Context(), &accountID); err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Account Deleted"})
}

func (h accountHandlerImpl) AccountPutHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	var accountInfo account.Account
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if !validateInput(&accountID, &accountInfo) {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	accountInfo.Account_id = accountID

	accountResponse, err := h.accountService.Update(c.Request().Context(), &accountInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, accountResponse)
}

func (h accountHandlerImpl) AccountGetReportHandler(c echo.Context) error {
	accountID := c.Param("account_id")
	if !validateInput(&accountID, nil) {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	accountReport, err := h.accountService.Report(c.Request().Context(), &accountID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusOK, accountReport)
}

func validateInput(param *string, accountRequest *account.Account) bool {
	if param != nil {
		if len(*param) > 20 {
			return false
		}
	}
	if accountRequest != nil {
		if len(accountRequest.Client_id) > 20 || len(accountRequest.Register_date) > 20 || len(accountRequest.User_id) > 20 {
			return false
		}
		if accountRequest.Agency_id <= 0 || accountRequest.User_id == "" {
			return false
		}
	}
	if accountRequest == nil && param == nil {
		return false
	}
	return true
}
