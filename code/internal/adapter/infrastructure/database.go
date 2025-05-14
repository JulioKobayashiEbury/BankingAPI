package adapter

import "BankingAPI/code/internal/domain"

func DeleteUserDB(userID int32) (int32, error) {
	id := int32(0)
	return id, nil
}

func CreateUserDB(user *domain.User) (*domain.User, error) {
	return &domain.User{}, nil
}

func UpdateUserDB(user *domain.User) (*domain.User, error) {
	// consult user
	// modify user
	// save user
	return &domain.User{}, nil
}
