package controller

import (
	"errors"
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type TransferHandler interface {
	TransferPostHandler(c echo.Context) error
	TransferGetHandler(c echo.Context) error
	TransferDeleteHandler(c echo.Context) error
}

type transferHandlerImpl struct {
	transferService service.TransferService
	accountService  service.AccountService
}

func NewTransferHandler(transferServe service.TransferService, accountServe service.AccountService) TransferHandler {
	return transferHandlerImpl{
		transferService: transferServe,
		accountService:  accountServe,
	}
}

func AddTransferEndPoints(server *echo.Echo, h TransferHandler) {
	server.POST("/transfers", h.TransferPostHandler)
	server.GET("/transfers/:transfer_id", h.TransferGetHandler)
	server.DELETE("/transfers/:transfer_id", h.TransferDeleteHandler)
}

func (h transferHandlerImpl) TransferPostHandler(c echo.Context) error {
	var newTransferInfo transfer.Transfer
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if _, err := h.authorizationForTransferEndpoints(&c, &newTransferInfo.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	transferResponse, err := h.transferService.ProcessNewTransfer(&newTransferInfo)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *transferResponse)
}

func (h transferHandlerImpl) TransferGetHandler(c echo.Context) error {
	transferID := c.Param("transfer_id")
	transferResponse, err := h.transferService.Get(&transferID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if _, err := h.authorizationForTransferEndpoints(&c, &transferResponse.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, *transferResponse)
}

func (h transferHandlerImpl) TransferDeleteHandler(c echo.Context) error {
	transferID := c.Param("transfer_id")
	transferResponse, err := h.transferService.Get(&transferID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if _, err := h.authorizationForTransferEndpoints(&c, &transferResponse.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	if err := h.transferService.Delete(&transferID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusOK, *transferResponse)
}

func (h transferHandlerImpl) authorizationForTransferEndpoints(c *echo.Context, accountID *string) (*string, *model.Erro) {
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
		return nil, &model.Erro{Err: errors.New("no match for user id"), HttpCode: http.StatusForbidden}
	}

	return nil, nil
}
