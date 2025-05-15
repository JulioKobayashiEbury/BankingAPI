package controller

import (
	"net/http"

	"BankingAPI/code/internal/domain"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func AddAccountEndPoints(server *echo.Echo) {
	server.GET("/accounts", AccountGetOrderFilterHandler)
	server.GET("/accounts/:account_id", AccountGetHandler)
	server.POST("/accounts", AccountPostHandler)
	server.DELETE("/accounts/:account_id", AccountDeleteHandler)
	server.PUT("/accounts/:account_id/withdrawal", AccountPutWithDrawalHandler)
	server.PUT("/accounts/:account_id/deposit", AccountPutDepositHandler)
}

func AccountPostHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountGetHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountDeleteHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutDepositHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutBlockHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutUnBlockHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPutAutomaticDebit(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}

func AccountPostTranferHandler(c echo.Context) error {
	var accountInfo domain.AccountRequest
	if err := c.Bind(&accountInfo); err != nil {
		log.Error().Msg(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	// talk to service

	return c.JSON(http.StatusOK, accountInfo)
}
