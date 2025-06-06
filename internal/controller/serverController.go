package controller

/* zerolog.SetGlobalLevel(zerolog.InfoLevel)
log.Info().Msg("Method not allowed")
w.WriteHeader(http.StatusMethodNotAllowed)
response := map[string]string{"error": "Method not allowed"}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
return
*/

import (
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/model/withdrawal"
	"BankingAPI/internal/service"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	documentLenghtIdeal = 14
	maxNameLenght       = 30
)

var (
	Repositories   model.RepositoryList
	Services       service.ServicesList
	DatabaseClient *firestore.Client
)

func Server() {
	server := echo.New()
	AddAccountEndPoints(server)
	AddClientsEndPoints(server)
	AddUsersEndPoints(server)

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if err := server.Start("localhost:25565"); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg("Server started on port 25565")
}

func InstantiateRepo() {
	Repositories = model.RepositoryList{
		UserDatabase:           user.NewUserFireStore(DatabaseClient),
		ClientDatabase:         client.NewClientFirestore(DatabaseClient),
		AccountDatabase:        account.NewAccountFirestore(DatabaseClient),
		AutomaticDebitDatabase: automaticdebit.NewAutoDebitFirestore(DatabaseClient),
		DepositDatabase:        deposit.NewDepositFirestore(DatabaseClient),
		TransferDatabase:       transfer.NewTransferFirestore(DatabaseClient),
		WithdrawalDatabase:     withdrawal.NewWithdrawalFirestore(DatabaseClient),
	}
}

func InstantiateServices() {
	getFilteredServe := service.NewGetFilteredService(
		Repositories.ClientDatabase,
		Repositories.AccountDatabase,
		Repositories.TransferDatabase,
		Repositories.DepositDatabase,
		Repositories.WithdrawalDatabase,
		Repositories.AutomaticDebitDatabase,
	)
	userServe := service.NewUserService(Repositories.UserDatabase, getFilteredServe)
	clientServe := service.NewClientService(Repositories.ClientDatabase, userServe, getFilteredServe)
	accountServe := service.NewAccountService(Repositories.AccountDatabase, userServe, clientServe, getFilteredServe)
	withdrawalServe := service.NewWithdrawalService(Repositories.WithdrawalDatabase, accountServe)
	depositServe := service.NewDepositService(Repositories.DepositDatabase, accountServe)
	automaticdebitServe := service.NewAutoDebit(Repositories.AutomaticDebitDatabase, withdrawalServe)
	transferServe := service.NewTransferService(Repositories.TransferDatabase, accountServe)

	Services = service.ServicesList{
		UserService:           userServe,
		ClientService:         clientServe,
		AccountService:        accountServe,
		WithdrawalService:     withdrawalServe,
		DepositService:        depositServe,
		AutomaticdebitService: automaticdebitServe,
		TransferService:       transferServe,
		GetFilteredService:    getFilteredServe,
	}
}
