package ports

import (
	"BankingAPI/code/internal/adapter"

	"github.com/labstack/echo/v4"
)

func UserPostHandler(c echo.Context) error {
	return c.JSON(adapter.UserPostAdapter(&c))
}

func UserPutBlockHandler(c echo.Context) error {
	return c.JSON(adapter.UserPutBlock(&c))
}

func UserPutUnBlockHandler(c echo.Context) error {
	return c.JSON(adapter.UserPutUnblock(&c))
}

func UserPutHandler(c echo.Context) error {
	return c.JSON(adapter.UserPutAdapter(&c))
}

func UserGetHandler(c echo.Context) error {
	return c.JSON(adapter.UserGetAdapter(&c))
}

func UserDeleteHandler(c echo.Context) error {
	return c.JSON(adapter.UserDeleteAdapter(&c))
}
