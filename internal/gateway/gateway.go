package gateway

import (
	"github.com/labstack/echo/v4"
)

type Gateway interface {
	Send(interface{}) *echo.HTTPError
}

type GatewaysList struct {
	ExternalTransferGateway Gateway
}
