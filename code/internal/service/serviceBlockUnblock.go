package service

import (
	controller "BankingAPI/code/internal/controller/objects"
)

func AccountBlock(account uint32) error {
	//get account from db
	//update account
	//put account into db again
	return nil
}

func AccountUnBlock(account uint32) error {
	return nil
}

func ClientBlock(client *controller.ClientRequest) error {
	return nil
}

func ClientUnBlock(client *controller.ClientRequest) error {
	return nil
}

func UserBlock(user *controller.UserRequest) error {
	return nil
}

func UserUnBlock(user *controller.UserRequest) error {
	return nil
}
