package adapter

import (
	"BankingAPI/code/internal/domain"
)

func DeleteUserDB(userID uint32) error {
	return nil
}

func CreateUserDB(user *domain.User) error {
	return nil
}

func UpdateUserDB(user *domain.User) error {
	// consult user
	// modify user
	// save user
	return nil
}

func CreateClientDB(client *domain.Client) error {
	return nil
}

func GetClientDB(client *domain.Client) error {
	return nil
}

func DeleteClientDB(clientID uint32) error {
	return nil
}

func UpdateClientDB(client *domain.Client) error {
	return nil
}

func CreateAccountDB(account *domain.Account) error {
	return nil
}

func GetAccountDB(account *domain.Account) error {
	return nil
}

func UpdateAccountDB(account *domain.Account) error {
	return nil
}

func DeleteAccountDB(accountID uint32) error {
	return nil
}

func SearchAccountByFilterOrder(filter string, order string, listOdAccounts *[]domain.Account) error {
	return nil
}
