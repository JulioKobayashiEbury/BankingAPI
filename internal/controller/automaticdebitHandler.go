package controller

import (
	"net/http"

	"BankingAPI/internal/model"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/service"

	"github.com/labstack/echo/v4"
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

func AddAutodebitEndPoints(group *echo.Group, h AutodebitHandler) {
	group.POST("/autodebits", h.AutodebitPostHandler)
	group.GET("/autodebits/:debit_id", h.AutodebitGetHandler)
	group.DELETE("/autodebits/:debit_id", h.AutodebitDeleteHandler)
}

func (h autodebitHandlerImpl) AutodebitPostHandler(c echo.Context) error {
	var newAutoDebit automaticdebit.AutomaticDebit
	if err := c.Bind(&newAutoDebit); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	autodebitResponse, err := h.automaticdebitService.ProcessNewAutomaticDebit(c.Request().Context(), &newAutoDebit)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusAccepted, *autodebitResponse)
}

func (h autodebitHandlerImpl) AutodebitGetHandler(c echo.Context) error {
	autodebitID := c.Param("debit_id")
	autodebitResponse, err := h.automaticdebitService.Get(c.Request().Context(), &autodebitID)
	if err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusOK, *autodebitResponse)
}

func (h autodebitHandlerImpl) AutodebitDeleteHandler(c echo.Context) error {
	autodebitID := c.Param("debit_id")

	if err := h.automaticdebitService.Delete(c.Request().Context(), &autodebitID); err != nil {
		return c.JSON(err.Code, err.Error())
	}

	return c.JSON(http.StatusNoContent, model.StandartResponse{Message: "Automatic debit deleted successfully"})
}
