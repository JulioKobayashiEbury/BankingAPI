package service

import (
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
	Authenticate(typeID *string, password *string) (bool, *model.Erro)
	GenerateToken(typeID *string) (*string, *model.Erro)
}

type UserService interface {
	Create(*user.User) (*user.User, *model.Erro)
	Delete(*string) *model.Erro
	Get(*string) (*user.User, *model.Erro)
	Update(*user.User) (*user.User, *model.Erro)
	GetAll() (*[]user.User, *model.Erro)
	Report(*string) (*user.UserReport, *model.Erro)
}

type ClientService interface {
	Create(*client.Client) (*client.Client, *model.Erro)
	Delete(*string) *model.Erro
	Get(*string) (*client.Client, *model.Erro)
	Update(*client.Client) (*client.Client, *model.Erro)
	GetAll() (*[]client.Client, *model.Erro)
	Report(*string) (*client.ClientReport, *model.Erro)
}

type AccountService interface {
	Create(*account.Account) (*account.Account, *model.Erro)
	Delete(*string) *model.Erro
	Get(*string) (*account.Account, *model.Erro)
	Update(*account.Account) (*account.Account, *model.Erro)
	GetAll() (*[]account.Account, *model.Erro)
	Report(*string) (*account.AccountReport, *model.Erro)
}

type WithdrawalService interface {
	Create(*withdrawal.Withdrawal) (*withdrawal.Withdrawal, *model.Erro)
	Delete(*string) *model.Erro
	Get(*string) (*withdrawal.Withdrawal, *model.Erro)
	GetAll() (*[]withdrawal.Withdrawal, *model.Erro)
	ProcessWithdrawal(withdrawalRequest *withdrawal.Withdrawal) (*withdrawal.Withdrawal, *model.Erro)
}

type DepositService interface {
	Create(*deposit.Deposit) (*deposit.Deposit, *model.Erro)
	Delete(*string) *model.Erro
	Get(*string) (*deposit.Deposit, *model.Erro)
	GetAll() (*[]deposit.Deposit, *model.Erro)
	ProcessDeposit(depositRequest *deposit.Deposit) (*deposit.Deposit, *model.Erro)
}

type AutomaticDebitService interface {
	Create(*automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro)
	Delete(*string) *model.Erro
	Get(*string) (*automaticdebit.AutomaticDebit, *model.Erro)
	GetAll() (*[]automaticdebit.AutomaticDebit, *model.Erro)
	ProcessNewAutomaticDebit(autoDebit *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *model.Erro)
	CheckAutomaticDebits()
	Scheduled()
}

type TransferService interface {
	Create(*transfer.Transfer) (*transfer.Transfer, *model.Erro)
	Delete(*string) *model.Erro
	Get(id *string) (*transfer.Transfer, *model.Erro)
	GetAll() (*[]transfer.Transfer, *model.Erro)
	ProcessNewTransfer(*transfer.Transfer) (*transfer.Transfer, *model.Erro)
}
