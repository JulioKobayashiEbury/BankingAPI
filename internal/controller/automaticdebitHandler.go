package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

type AutodebitHandler interface {
	AutodebitPostHandler(c echo.Context) error
	AutodebitGetHandler(c echo.Context) error
	AutodebitDeleteHandler(c echo.Context) error
}

type autodebitHandlerImpl struct {
	automaticdebitService service.AutomaticDebitService
	accountService        service.AccountService
}

func NewAutodebitHandler(automaticdebitServe service.AutomaticDebitService, accountServe service.AccountService) AutodebitHandler {
	return autodebitHandlerImpl{
		automaticdebitService: automaticdebitServe,
		accountService:        accountServe,
	}
}

func AddAutodebitEndPoints(server *echo.Echo, h AutodebitHandler) {
	server.POST("/autodebits", h.AutodebitPostHandler)
	server.GET("/autodebits/:debit_id", h.AutodebitGetHandler)
	server.DELETE("/autodebits/:debit_id", h.AutodebitDeleteHandler)
}

func (h autodebitHandlerImpl) AutodebitPostHandler(c echo.Context) error {
	var newAutoDebit automaticdebit.AutomaticDebit
	if err := c.Bind(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	autodebitResponse, err := h.automaticdebitService.ProcessNewAutomaticDebit(&newAutoDebit)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusAccepted, *autodebitResponse)
}

func (h autodebitHandlerImpl) AutodebitGetHandler(c echo.Context) error {
	autodebitID := c.Param("debit_id")
	autodebitResponse, err := h.automaticdebitService.Get(&autodebitID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *autodebitResponse)
}

func (h autodebitHandlerImpl) AutodebitDeleteHandler(c echo.Context) error {
	autodebitID := c.Param("debit_id")

	if err := h.automaticdebitService.Delete(&autodebitID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusNoContent, model.StandartResponse{Message: "Automatic debit deleted successfully"})
}
