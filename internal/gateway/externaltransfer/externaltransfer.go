package externaltransfer

import (
	"BankingAPI/internal/gateway"
	"BankingAPI/internal/model"
)

type externalTransferImpl struct{}

func NewExternalTransferGateway() gateway.Gateway {
	return externalTransferImpl{}
}

func (ex externalTransferImpl) Send(interface{}) *model.Erro { // money is leaving this system (inside -> outside)
	// do nothing, interact with partner's bank API
	return nil
}
