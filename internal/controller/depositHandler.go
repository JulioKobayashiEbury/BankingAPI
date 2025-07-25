package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
)

type DepositHandler interface {
	PostDepositHandler(c echo.Context) error
	DeleteDepositHandler(c echo.Context) error
	GetDepositHandler(c echo.Context) error
}

type depositHandlerImpl struct {
	depositService service.DepositService
	accountService service.AccountService
}

func NewDeposithandler(depositServe service.DepositService, accountServe service.AccountService) DepositHandler {
	return depositHandlerImpl{
		depositService: depositServe,
		accountService: accountServe,
	}
}

func AddDepositsEndPoints(group *echo.Group, h DepositHandler) {
	group.POST("/deposits", h.PostDepositHandler)
	group.GET("/deposits/:deposit_id", h.GetDepositHandler)
	group.DELETE("/deposits/:deposit_id", h.DeleteDepositHandler)
}

func (h depositHandlerImpl) PostDepositHandler(c echo.Context) error {
	var newDepositInfo deposit.Deposit
	if err := c.Bind(&newDepositInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if newDepositInfo.User_id == "" || newDepositInfo.Account_id == "" || newDepositInfo.Deposit <= 0 {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	depositResponse, err := h.depositService.ProcessDeposit(c.Request().Context(), &newDepositInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusCreated, *depositResponse)
}

func (h depositHandlerImpl) DeleteDepositHandler(c echo.Context) error {
	depositID := c.Param("deposit_id")

	if err := h.depositService.Delete(c.Request().Context(), &depositID); err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Deposit deleted succesfully!"})
}

func (h depositHandlerImpl) GetDepositHandler(c echo.Context) error {
	depositID := c.Param("deposit_id")

	depositInfo, err := h.depositService.Get(c.Request().Context(), &depositID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, *depositInfo)
}
