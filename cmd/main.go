package main

import (
	"BankingAPI/internal/controller"
	"os"
)

func init() {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "0.0.0.0:9000")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "banking")
}

func main() {
	controller.Server()
}
