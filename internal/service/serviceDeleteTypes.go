package service

import "BankingAPI/internal/model"

func AccountDelete(accountID string) error {
	if err := model.DeleteObject(&accountID, model.AccountsPath); err != nil {
		return err
	}
	return nil
}

func ClientDelete(clientID string) error {
	if err := model.DeleteObject(&clientID, model.ClientPath); err != nil {
		return err
	}
	return nil
}

func UserDelete(userID string) error {
	if err := model.DeleteObject(&userID, model.UsersPath); err != nil {
		return err
	}
	return nil
}
