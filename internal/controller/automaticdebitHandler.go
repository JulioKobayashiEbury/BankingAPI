package controller

import (
	"errors"
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

	if _, err := h.authorizationForAutodebitEndpoints(&c, &newAutoDebit.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	autodebitResponse, err := h.automaticdebitService.ProcessNewAutomaticDebit(&newAutoDebit)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusAccepted, *autodebitResponse)
}

func (h autodebitHandlerImpl) AutodebitGetHandler(c echo.Context) error {
	autodebitID := c.Param("debit_id")
	if _, err := h.authorizationForAutodebitEndpoints(&c, &autodebitID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	autodebitResponse, err := h.automaticdebitService.Get(&autodebitID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if _, err := h.authorizationForAutodebitEndpoints(&c, &autodebitResponse.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusOK, *autodebitResponse)
}

func (h autodebitHandlerImpl) AutodebitDeleteHandler(c echo.Context) error {
	autodebitID := c.Param("debit_id")
	autodebitResponse, err := h.automaticdebitService.Get(&autodebitID)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if _, err := h.authorizationForAutodebitEndpoints(&c, &autodebitResponse.Account_id); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	if err := h.automaticdebitService.Delete(&autodebitID); err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}

	return c.JSON(http.StatusNoContent, model.StandartResponse{Message: "Automatic debit deleted successfully"})
}

func (h autodebitHandlerImpl) authorizationForAutodebitEndpoints(c *echo.Context, accountID *string) (*string, *model.Erro) {
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
