package application

import (
	"errors"

	"BankingAPI/code/internal/domain"

	"github.com/rs/zerolog/log"
)

func Deposit(account *domain.Account, amount float64) error {
	if amount <= 0 {
		return errors.New("Deposit amount must be positive")
	}
	account.AddBalance(amount)
	log.Info().Msgf("Deposited %.2f to account %d", amount, account.GetAccountId())
	return nil
}

func Withdraw(account *domain.Account, amount float64) error {
	if amount <= 0 {
		return errors.New("Withdrawal amount must be positive")
	}
	if account.GetBalance() < amount {
		return errors.New("Insufficient funds")
	}
	account.SubtractBalance(amount)
	log.Info().Msgf("Withdrew %.2f from account %d", amount, account.GetAccountId())
	return nil
}

func Transfer(sourceAccount *domain.Account, targetAccount *domain.Account, amount float64) error {
	if amount <= 0 {
		return errors.New("Transfer amount must be positive")
	}
	if sourceAccount.GetBalance() < amount {
		return errors.New("Insufficient funds")
	}
	sourceAccount.TransferBalance(amount, targetAccount)
	log.Info().Msgf("Transferred %.2f from account %d to account %d", amount, sourceAccount.GetAccountId(), targetAccount.GetAccountId())
	return nil
}

func GetBalance(account *domain.Account) float64 {
	return account.GetBalance()
}
