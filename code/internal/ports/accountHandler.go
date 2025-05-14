package ports

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccountPostHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountGetHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountGetOrderFilterHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountDeleteHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutDepositHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func AccountPutWithDrawalHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
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
