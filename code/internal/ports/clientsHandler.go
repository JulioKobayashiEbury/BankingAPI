package ports

import (
	"BankingAPI/code/internal/adapter"

	"github.com/labstack/echo/v4"
)

func ClientPostHandler(c echo.Context) error {
	return c.JSON(adapter.ClientPostAdapter(&c))
}

func ClientGetHandler(c echo.Context) error {
	return c.JSON(adapter.ClientGetAdapter(&c))
}

func ClientDeleteHandler(c echo.Context) error {
	return c.JSON(adapter.ClientDeleteAdapter(&c))
}

func ClientPutHandler(c echo.Context) error {
	return c.JSON(adapter.ClientPutAdapter(&c))
}

func ClientPutBlockHandler(c echo.Context) error {
	return c.JSON(adapter.ClientPutBlockAdapter(&c))
}

func ClientPutUnBlockHandler(c echo.Context) error {
	return c.JSON(adapter.ClientPutUnBlockAdapter(&c))
}
