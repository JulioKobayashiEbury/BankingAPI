package gateway

import (
	"github.com/labstack/echo"
)

type Gateway interface {
	Send(interface{}) *echo.HTTPError
}

type GatewaysList struct {
	ExternalTransferGateway Gateway
}
