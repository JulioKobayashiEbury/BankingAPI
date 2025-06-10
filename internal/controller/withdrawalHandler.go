package controller

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/withdrawal"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
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

func AddWithdrawalEndPoints(server *echo.Echo, h WithdrawalHandler) {
	server.POST("/withdrawals", h.PostWithdrawalHandler)
	server.GET("/withdrawals/:withdrawal_id", h.GetWithdrawalHandler)
	server.DELETE("withdrawals/:withdrawal_id", h.DeleteWithdrawalHandler)
}

func (h withdrawalHandlerImpl) PostWithdrawalHandler(c echo.Context) error {
	var withdrawalInfo withdrawal.Withdrawal
	if err := c.Bind(&withdrawalInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := h.authorizationForWithdrawalEndPoints(&c, &withdrawalInfo.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	withdrawalReponse, err := h.withdrawalService.ProcessWithdrawal(&withdrawalInfo)
	if err != nil {
		c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusCreated, *withdrawalReponse)
}

func (h withdrawalHandlerImpl) DeleteWithdrawalHandler(c echo.Context) error {
	withdrawalID := c.Param("withdrawal_id")

	withdrawal, err := h.withdrawalService.Get(&withdrawalID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.authorizationForWithdrawalEndPoints(&c, &withdrawal.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.withdrawalService.Delete(&withdrawalID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Withdrawal Deleted!"})
}

func (h withdrawalHandlerImpl) GetWithdrawalHandler(c echo.Context) error {
	withdrawalID := c.Param("withdrawal_id")

	withdrawal, err := h.withdrawalService.Get(&withdrawalID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.authorizationForWithdrawalEndPoints(&c, &withdrawal.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	withdrawalResponse, err := h.withdrawalService.Get(&withdrawalID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *withdrawalResponse)
}

func (h withdrawalHandlerImpl) authorizationForWithdrawalEndPoints(c *echo.Context, accountID *string) *model.Erro {
	authorizationHeader := (*c).Request().Header.Get((echo.HeaderAuthorization))

	claims, err := service.Authorize(&authorizationHeader)
	if err != nil {
		return err
	}

	account, err := h.accountService.Get(accountID)
	if err != nil {
		return err
	}
	if account.User_id != claims.Id {
		log.Error().Msg("User ID does not match with accounts User ID")
		return &model.Erro{Err: errors.New("No match for user id"), HttpCode: http.StatusForbidden}
	}

	return nil
}
