package main

import (
	"BankingAPI/code/internal/ports"
)

var (
	UserIDCounter    int = 0
	ClientIDCounter  int = 0
	AccountIDCounter int = 0
)

func main() {
	ports.Server()
}
