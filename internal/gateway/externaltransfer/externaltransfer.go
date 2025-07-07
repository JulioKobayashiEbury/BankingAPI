package externaltransfer

import (
	"BankingAPI/internal/gateway"

	"github.com/labstack/echo"
)

type externalTransferImpl struct{}

func NewExternalTransferGateway() gateway.Gateway {
	return externalTransferImpl{}
}

func (ex externalTransferImpl) Send(interface{}) *echo.HTTPError { // money is leaving this system (inside -> outside)
	// do nothing, interact with partner's bank API
	return nil
}
