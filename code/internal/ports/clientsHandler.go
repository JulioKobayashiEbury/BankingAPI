package ports

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ClientPostHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func ClientGetHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func ClientDeleteHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func ClientPutHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func ClientPutBlockHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func ClientPutUnBlockHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
