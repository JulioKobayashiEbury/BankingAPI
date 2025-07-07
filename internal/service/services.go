package service

import (
	"context"

	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/model/withdrawal"

	"github.com/labstack/echo"
)

type ServicesList struct {
	UserService           UserService
	ClientService         ClientService
	AccountService        AccountService
	WithdrawalService     WithdrawalService
	DepositService        DepositService
	AutomaticdebitService AutomaticDebitService
	TransferService       TransferService
	AuthenticationService Authentication
}

type Authentication interface {
	Authenticate(ctx context.Context, typeID *string, password *string) (bool, *echo.HTTPError)
	GenerateToken(ctx context.Context, typeID *string) (*string, *echo.HTTPError)
}

type UserService interface {
	Create(context.Context, *user.User) (*user.User, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*user.User, *echo.HTTPError)
	Update(context.Context, *user.User) (*user.User, *echo.HTTPError)
	GetAll(context.Context) (*[]user.User, *echo.HTTPError)
	Report(context.Context, *string) (*user.UserReport, *echo.HTTPError)
}

type ClientService interface {
	Create(context.Context, *client.Client) (*client.Client, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*client.Client, *echo.HTTPError)
	Update(context.Context, *client.Client) (*client.Client, *echo.HTTPError)
	GetAll(context.Context) (*[]client.Client, *echo.HTTPError)
	Report(context.Context, *string) (*client.ClientReport, *echo.HTTPError)
}

type AccountService interface {
	Create(context.Context, *account.Account) (*account.Account, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*account.Account, *echo.HTTPError)
	Update(context.Context, *account.Account) (*account.Account, *echo.HTTPError)
	GetAll(context.Context) (*[]account.Account, *echo.HTTPError)
	Report(context.Context, *string) (*account.AccountReport, *echo.HTTPError)
}

type WithdrawalService interface {
	Create(context.Context, *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*withdrawal.Withdrawal, *echo.HTTPError)
	GetAll(context.Context) (*[]withdrawal.Withdrawal, *echo.HTTPError)
	ProcessWithdrawal(context.Context, *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *echo.HTTPError)
}

type DepositService interface {
	Create(context.Context, *deposit.Deposit) (*deposit.Deposit, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*deposit.Deposit, *echo.HTTPError)
	GetAll(context.Context) (*[]deposit.Deposit, *echo.HTTPError)
	ProcessDeposit(context.Context, *deposit.Deposit) (*deposit.Deposit, *echo.HTTPError)
}

type AutomaticDebitService interface {
	Create(context.Context, *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*automaticdebit.AutomaticDebit, *echo.HTTPError)
	GetAll(context.Context) (*[]automaticdebit.AutomaticDebit, *echo.HTTPError)
	ProcessNewAutomaticDebit(context.Context, *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *echo.HTTPError)
	CheckAutomaticDebits()
	Scheduled()
}

type TransferService interface {
	Create(context.Context, *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*transfer.Transfer, *echo.HTTPError)
	GetAll(context.Context) (*[]transfer.Transfer, *echo.HTTPError)
	ProcessNewTransfer(context.Context, *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError)
	ProcessExternalTransfer(context.Context, *transfer.Transfer) (*transfer.Transfer, *echo.HTTPError)
}
