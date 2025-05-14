package adapter

import (
	"BankingAPI/code/internal/domain"
)

func DeleteUserDB(userID int32) (int32, error) {
	id := int32(0)
	return id, nil
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

func DeleteClientDB(clientID int32) (int32, error) {
	return clientID, nil
}

func UpdateClientDB(client *domain.Client) error {
	return nil
}
