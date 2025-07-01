package service

import (
	"context"
	"errors"
	"net/http"

	"BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/client"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/model/withdrawal"
)

var ErrorMissingCredentials = &model.Erro{Err: errors.New("missing credentials"), HttpCode: http.StatusBadRequest}

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
	Authenticate(ctx context.Context, typeID *string, password *string) (bool, *model.Erro)
	GenerateToken(ctx context.Context, typeID *string) (*string, *model.Erro)
}

type UserService interface {
	Create(context.Context, *user.User) (*user.User, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*user.User, *model.Erro)
	Update(context.Context, *user.User) (*user.User, *model.Erro)
	GetAll(context.Context) (*[]user.User, *model.Erro)
	Report(context.Context, *string) (*user.UserReport, *model.Erro)
}

type ClientService interface {
	Create(context.Context, *client.Client) (*client.Client, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*client.Client, *model.Erro)
	Update(context.Context, *client.Client) (*client.Client, *model.Erro)
	GetAll(context.Context) (*[]client.Client, *model.Erro)
	Report(context.Context, *string) (*client.ClientReport, *model.Erro)
}

type AccountService interface {
	Create(context.Context, *account.Account) (*account.Account, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*account.Account, *model.Erro)
	Update(context.Context, *account.Account) (*account.Account, *model.Erro)
	GetAll(context.Context) (*[]account.Account, *model.Erro)
	Report(context.Context, *string) (*account.AccountReport, *model.Erro)
}

type WithdrawalService interface {
	Create(context.Context, *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*withdrawal.Withdrawal, *model.Erro)
	GetAll(context.Context) (*[]withdrawal.Withdrawal, *model.Erro)
	ProcessWithdrawal(context.Context, *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *model.Erro)
}

type DepositService interface {
	Create(context.Context, *deposit.Deposit) (*deposit.Deposit, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*deposit.Deposit, *model.Erro)
	GetAll(context.Context) (*[]deposit.Deposit, *model.Erro)
	ProcessDeposit(context.Context, *deposit.Deposit) (*deposit.Deposit, *model.Erro)
}

type AutomaticDebitService interface {
	Create(context.Context, *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*automaticdebit.AutomaticDebit, *model.Erro)
	GetAll(context.Context) (*[]automaticdebit.AutomaticDebit, *model.Erro)
	ProcessNewAutomaticDebit(context.Context, *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro)
	CheckAutomaticDebits()
	Scheduled()
}

type TransferService interface {
	Create(context.Context, *transfer.Transfer) (*transfer.Transfer, *model.Erro)
	Delete(context.Context, *string) *model.Erro
	Get(context.Context, *string) (*transfer.Transfer, *model.Erro)
	GetAll(context.Context) (*[]transfer.Transfer, *model.Erro)
	ProcessNewTransfer(context.Context, *transfer.Transfer) (*transfer.Transfer, *model.Erro)
	ProcessExternalTransfer(context.Context, *transfer.Transfer) (*transfer.Transfer, *model.Erro)
}
