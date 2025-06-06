package deposit

import "BankingAPI/internal/model"

var singleton *model.RepositoryInterface

type MockDepositRepository struct {
	UserMap *map[string]Deposit
}
