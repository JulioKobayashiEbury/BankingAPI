package service

import "BankingAPI/internal/repository"

func AccountDelete(accountID uint32) error {
	if err := repository.DeleteObject(accountID, repository.AccountsPath); err != nil {
		return err
	}
	return nil
}

func ClientDelete(clientID uint32) error {
	if err := repository.DeleteObject(clientID, repository.ClientPath); err != nil {
		return err
	}
	return nil
}

func UserDelete(userID uint32) error {
	if err := repository.DeleteObject(userID, repository.UsersPath); err != nil {
		return err
	}
	return nil
}
