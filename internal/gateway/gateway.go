package gateway

import (
	"BankingAPI/internal/model"
)

type Gateway interface {
	Send(interface{}) *model.Erro
}

type GatewaysList struct {
	ExternalTransferGateway Gateway
}
