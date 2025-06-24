package externaltransfer

import (
	"net/http"

	"BankingAPI/internal/gateway"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/service"

	"github.com/labstack/echo"
)

type externalTransferImpl struct {
	transfer service.TransferService
}

func AddExternalTransferEndpoint(server *echo.Echo, externalTransferGateway gateway.Gateway) {
	server.POST("/external-transfer", externalTransferGateway.Receive) // transfer coming from outside and entering the bank that is using this API (outside -> inside)
}

func (ex externalTransferImpl) Send(interface{}) error { // money is leaving this system (inside -> outside)
	// do nothing, interact with partner's bank API
	return nil
}

func (ex externalTransferImpl) Receive(c echo.Context) error {
	var transferRequest transfer.Transfer
	if err := c.Bind(&transferRequest); err != nil {
		httpErr := err.(*echo.HTTPError)
		return c.JSON(httpErr.Code, httpErr.Internal.Error())
	}
	// Checar cliente... conta... compliance...
	createdTransfer, err := ex.transfer.Create(&transferRequest)
	if err != nil {
		return c.JSON(err.HttpCode, err.Err.Error())
	}
	return c.JSON(http.StatusAccepted, *createdTransfer)
}
