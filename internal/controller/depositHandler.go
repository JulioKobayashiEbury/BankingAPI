package controller

import (
	"errors"
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
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

func AddDepositsEndPoints(server *echo.Echo, h DepositHandler) {
	server.POST("/deposits", h.PostDepositHandler)
	server.GET("/deposits/:deposit_id", h.GetDepositHandler)
	server.DELETE("/deposits/:deposit_id", h.DeleteDepositHandler)
}

func (h depositHandlerImpl) PostDepositHandler(c echo.Context) error {
	var newDepositInfo deposit.Deposit
	if err := c.Bind(&newDepositInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := h.authorizationForDepositsEndPoints(&c, &newDepositInfo.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	depositResponse, err := h.depositService.ProcessDeposit(&newDepositInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusCreated, *depositResponse)
}

func (h depositHandlerImpl) DeleteDepositHandler(c echo.Context) error {
	depositID := c.Param("deposit_id")

	depositInfo, err := h.depositService.Get(&depositID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.authorizationForDepositsEndPoints(&c, &depositInfo.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.depositService.Delete(&depositID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Deposit deleted succesfully!"})
}

func (h depositHandlerImpl) GetDepositHandler(c echo.Context) error {
	depositID := c.Param("deposit_id")

	depositInfo, err := h.depositService.Get(&depositID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.authorizationForDepositsEndPoints(&c, &depositInfo.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *depositInfo)
}

func (h depositHandlerImpl) authorizationForDepositsEndPoints(c *echo.Context, accountID *string) *model.Erro {
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
		return &model.Erro{Err: errors.New("no match for user id"), HttpCode: http.StatusForbidden}
	}

	return nil
}
