package gateway

import "github.com/labstack/echo"

type Gateway interface {
	Send(interface{}) error
	Receive(c echo.Context) error
}
