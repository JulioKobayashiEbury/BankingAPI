package controller

import (
	"net/http"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type TransferHandler interface {
	TransferPostHandler(c echo.Context) error
	TransferGetHandler(c echo.Context) error
	TransferDeleteHandler(c echo.Context) error
	ExternalTransferPostHandler(c echo.Context) error
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

func AddTransferEndPoints(group *echo.Group, h TransferHandler) {
	group.POST("/transfers", h.TransferPostHandler)
	group.POST("/external-transfers", h.TransferPostHandler)
	group.GET("/transfers/:transfer_id", h.TransferGetHandler)
	group.DELETE("/transfers/:transfer_id", h.TransferDeleteHandler)
}

func (h transferHandlerImpl) TransferPostHandler(c echo.Context) error {
	var newTransferInfo transfer.Transfer
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if newTransferInfo.Account_id == "" || newTransferInfo.Account_to == "" || newTransferInfo.Value <= 0 || newTransferInfo.User_id == "" || newTransferInfo.User_to == "" {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	transferResponse, err := h.transferService.ProcessNewTransfer(c.Request().Context(), &newTransferInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, *transferResponse)
}

func (h transferHandlerImpl) TransferGetHandler(c echo.Context) error {
	transferID := c.Param("transfer_id")
	transferResponse, err := h.transferService.Get(c.Request().Context(), &transferID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusOK, *transferResponse)
}

func (h transferHandlerImpl) TransferDeleteHandler(c echo.Context) error {
	transferID := c.Param("transfer_id")
	if err := h.transferService.Delete(c.Request().Context(), &transferID); err != nil {
		return c.JSON(err.Code, err.Error())
	}
	return c.JSON(http.StatusOK, model.StandartResponse{Message: "Transfer Deleted!"})
}

func (h transferHandlerImpl) ExternalTransferPostHandler(c echo.Context) error {
	var newTransferInfo transfer.Transfer
	if err := c.Bind(&newTransferInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if newTransferInfo.Account_id == "" || newTransferInfo.Account_to == "" || newTransferInfo.Value <= 0 || newTransferInfo.User_id == "" || newTransferInfo.User_to == "" {
		return c.JSON(http.StatusBadRequest, model.StandartResponse{Message: "Parameters are not ideal"})
	}

	transferResponse, err := h.transferService.ProcessExternalTransfer(c.Request().Context(), &newTransferInfo)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, *transferResponse)
}
