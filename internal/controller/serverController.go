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
	_ "embed"
	"net/http"

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
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	middleware := middleware.NewUserAuthMiddleware(services.UserService)
	external := server.Group("/external")
	internal := server.Group("/internal")
	external.Use(echo.MiddlewareFunc(middleware.AuthorizeMiddleware))

	setupOpenApiDocs(internal)

	AddAuthenticationEndpoints(server, NewAuthenticationHandler(services.AuthenticationService))

	AddAccountEndPoints(external, NewAccountHandler(services.AccountService))
	AddClientsEndPoints(external, NewClientHandler(services.ClientService))
	AddUsersEndPoints(external, NewUserHandler(services.UserService, services.AuthenticationService))
	AddTransferEndPoints(external, NewTransferHandler(services.TransferService, services.AccountService))
	AddAutodebitEndPoints(external, NewAutodebitHandler(services.AutomaticdebitService, services.AccountService))
	AddDepositsEndPoints(external, NewDeposithandler(services.DepositService, services.AccountService))
	AddWithdrawalEndPoints(external, NewWithdrawalHandler(services.WithdrawalService, services.AccountService))

	services.AutomaticdebitService.Scheduled()

	if err := server.Start("localhost:25565"); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg("Server started on port 25565")
}

//go:embed docs/swagger.yaml
var swagger []byte

func setupOpenApiDocs(group *echo.Group) {
	group.GET("/swagger.yaml", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/yaml; charset=utf-8")
		c.Response().Header().Set(echo.HeaderContentDisposition, "inline")

		return c.Blob(http.StatusOK, "text/yaml; charset=utf-8", swagger)
	})

	group.GET("/docs/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{"/internal/swagger.yaml"}
	}))
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
