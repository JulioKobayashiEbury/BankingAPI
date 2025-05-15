package ports

import (
	"net/http"

	"BankingAPI/code/internal/adapter"

	"github.com/labstack/echo/v4"
)

func AccountPostHandler(c echo.Context) error {
	return c.JSON(adapter.AccountPostAdapter(&c))
}

func AccountGetHandler(c echo.Context) error {
	return c.JSON(adapter.AccountGetAdapter(&c))
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	return c.JSON(adapter.AccountGetByOrderFilterAdapter(&c))
}

func AccountDeleteHandler(c echo.Context) error {
	return c.JSON(adapter.AccountDeleteAdapter(&c))
}

func AccountPutHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutDepositHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	return c.JSON(adapter.AccountPutWithdrawalAdapter(&c))
}

func AccountPutBlockHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutUnBlockHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutAutomaticDebit(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPostTranferHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
