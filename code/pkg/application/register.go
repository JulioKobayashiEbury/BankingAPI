package application

import (
	"BankingAPI/code/pkg/domain"

	"github.com/rs/zerolog/log"
)

func UpdatePassword(account *domain.Account, newPassword string) {
	account.UpdatePassword(newPassword)
	log.Info().Msgf("Updated password for account %d", account.GetAccountId())
}

func ActivateAccount(account *domain.Account) {
	account.Activate()
	log.Info().Msgf("Activated account %d", account.GetAccountId())
}

func DeactivateAccount(account *domain.Account) {
	account.Deactivate()
	log.Info().Msgf("Deactivated account %d", account.GetAccountId())
}

func IsActive(account *domain.Account) bool {
	return account.IsActive()
}

func SetAccountStatus(account *domain.Account, status bool) {
	account.SetStatus(status)
	log.Info().Msgf("Set account %d status to %v", account.GetAccountId(), status)
}
