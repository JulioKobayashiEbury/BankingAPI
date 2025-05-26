package service

import (
	"fmt"

	model "BankingAPI/internal/model/types"
)

func GenerateReportByAccount(accountID *string) (*model.AccountReport, *model.Erro) {
	accountInfo, err := Account(*accountID)
	if err != nil {
		return nil, err
	}
	fmt.Print(accountInfo)
	return nil, nil
}
