package controller

import (
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/withdrawal"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
)

type WithdrawalHandler interface {
	PostWithdrawalHandler(c echo.Context) error
	DeleteWithdrawalHandler(c echo.Context) error
	GetWithdrawalHandler(c echo.Context) error
}

type withdrawalHandlerImpl struct {
	withdrawalService service.WithdrawalService
	accountService    service.AccountService
}

func NewWithdrawalHandler(withdrawalServe service.WithdrawalService, accontServe service.AccountService) WithdrawalHandler {
	return withdrawalHandlerImpl{
		withdrawalService: withdrawalServe,
		accountService:    accontServe,
	}
}

func AddWithdrawalEndPoints(group *echo.Group, h WithdrawalHandler) {
	group.POST("/withdrawals", h.PostWithdrawalHandler)
	group.GET("/withdrawals/:withdrawal_id", h.GetWithdrawalHandler)
	group.DELETE("withdrawals/:withdrawal_id", h.DeleteWithdrawalHandler)
}

func (h withdrawalHandlerImpl) PostWithdrawalHandler(c echo.Context) error {
	var withdrawalInfo withdrawal.Withdrawal
	if err := c.Bind(&withdrawalInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if withdrawalInfo.Account_id == "" || withdrawalInfo.Withdrawal <= 0 || withdrawalInfo.User_id == "" {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	withdrawalReponse, err := h.withdrawalService.ProcessWithdrawal(c.Request().Context(), &withdrawalInfo)
	if err != nil {
		c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusCreated, *withdrawalReponse)
}

func (h withdrawalHandlerImpl) DeleteWithdrawalHandler(c echo.Context) error {
	withdrawalID := c.Param("withdrawal_id")

	if err := h.withdrawalService.Delete(c.Request().Context(), &withdrawalID); err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Withdrawal Deleted!"})
}

func (h withdrawalHandlerImpl) GetWithdrawalHandler(c echo.Context) error {
	withdrawalID := c.Param("withdrawal_id")

	withdrawalResponse, err := h.withdrawalService.Get(c.Request().Context(), &withdrawalID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, *withdrawalResponse)
}
