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
	"BankingAPI/internal/gateway"
	"BankingAPI/internal/gateway/externaltransfer"
	"BankingAPI/internal/middleware"
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
	"github.com/rs/zerolog/log"
)

type RepositoryList struct {
	UserDatabase           user.UserRepository
	ClientDatabase         client.ClientRepository
	AccountDatabase        account.AccountRepository
	AutomaticDebitDatabase automaticdebit.AutoDebitRepository
	DepositDatabase        deposit.DepositRepository
	TransferDatabase       transfer.TransferRepository
	WithdrawalDatabase     withdrawal.WithdrawalRepository
}

const (
	documentLenghtForUser   = 14
	documentLenghtForClient = 11
	maxNameLenght           = 30
)

func Server(services *service.ServicesList) {
	server := echo.New()

	AddAccountEndPoints(server, NewAccountHandler(services.AccountService))
	AddClientsEndPoints(server, NewClientHandler(services.ClientService))
	AddUsersEndPoints(server, NewUserHandler(services.UserService, services.AuthenticationService))

	middleware := middleware.NewUserAuthMiddleware(services.UserService)
	server.Use(echo.MiddlewareFunc(middleware.AuthorizeMiddleware))

	AddAuthenticationEndpoints(server, NewAuthenticationHandler(services.AuthenticationService))
	AddTransferEndPoints(server, NewTransferHandler(services.TransferService, services.AccountService))
	AddAutodebitEndPoints(server, NewAutodebitHandler(services.AutomaticdebitService, services.AccountService))
	AddDepositsEndPoints(server, NewDeposithandler(services.DepositService, services.AccountService))
	AddWithdrawalEndPoints(server, NewWithdrawalHandler(services.WithdrawalService, services.AccountService))

	services.AutomaticdebitService.Scheduled()

	if err := server.Start("localhost:25565"); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg("Server started on port 25565")
}

func InstantiateRepo(databaseClient *firestore.Client) *RepositoryList {
	return &RepositoryList{
		UserDatabase:           user.NewUserFireStore(databaseClient),
		ClientDatabase:         client.NewClientFirestore(databaseClient),
		AccountDatabase:        account.NewAccountFirestore(databaseClient),
		AutomaticDebitDatabase: automaticdebit.NewAutoDebitFirestore(databaseClient),
		DepositDatabase:        deposit.NewDepositFirestore(databaseClient),
		TransferDatabase:       transfer.NewTransferFirestore(databaseClient),
		WithdrawalDatabase:     withdrawal.NewWithdrawalFirestore(databaseClient),
	}
}

func InstantiateServices(repositories *RepositoryList, gateways *gateway.GatewaysList) *service.ServicesList {
	userServe := service.NewUserService(repositories.UserDatabase, repositories.ClientDatabase)
	clientServe := service.NewClientService(repositories.ClientDatabase, userServe, repositories.AccountDatabase)
	accountServe := service.NewAccountService(repositories.AccountDatabase,
		userServe,
		clientServe,
		repositories.WithdrawalDatabase,
		repositories.DepositDatabase,
		repositories.TransferDatabase,
		repositories.AutomaticDebitDatabase,
	)
	withdrawalServe := service.NewWithdrawalService(repositories.WithdrawalDatabase, accountServe)
	depositServe := service.NewDepositService(repositories.DepositDatabase, accountServe)
	automaticdebitServe := service.NewAutoDebit(repositories.AutomaticDebitDatabase, withdrawalServe)

	transferServe := service.NewTransferService(repositories.TransferDatabase, accountServe, userServe, gateways.ExternalTransferGateway)
	authentication := service.NewAuth(repositories.UserDatabase)

	return &service.ServicesList{
		UserService:           userServe,
		ClientService:         clientServe,
		AccountService:        accountServe,
		WithdrawalService:     withdrawalServe,
		DepositService:        depositServe,
		AutomaticdebitService: automaticdebitServe,
		TransferService:       transferServe,
		AuthenticationService: authentication,
	}
}

func InstantiateGateways() *gateway.GatewaysList {
	return &gateway.GatewaysList{
		ExternalTransferGateway: externaltransfer.NewExternalTransferGateway(),
	}
}
