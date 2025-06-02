package service_test

import (
	"testing"

	"BankingAPI/internal/model/deposit"
)

func TestProcessDeposit(t *testing.T) {
	type TestDeposit []struct {
		test     string
		value    *deposit.DepositRequest
		expected *deposit.DepositResponse
	}
	tests := TestDeposit{
		{
			test: "First deposit",
			value: &deposit.DepositRequest{
				Account_id: "accID",
				Client_id:  "cllID",
				User_id:    "usrID",
				Agency_id:  1,
				Deposit:    1502.82,
			},
		},
		{
			test: "Second deposit",
			value: &deposit.DepositRequest{
				Account_id: "accID2",
				Client_id:  "cllID2",
				User_id:    "usrID2",
				Agency_id:  2,
				Deposit:    2.82,
			},
		},
	}
	for _, deposit := range tests {
		// teste aqui

		t.Run(deposit.test, func(t *testing.T) {
			// ou aqui
		})
	}
}
