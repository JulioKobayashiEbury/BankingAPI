package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/deposit"
	"BankingAPI/internal/model/transfer"
	"BankingAPI/internal/model/withdrawal"

	"github.com/rs/zerolog/log"
)

type accountServiceImpl struct {
	userService     UserService
	clientService   ClientService
	accountDatabase account.AccountRepository

	withdrawalDatabase withdrawal.WithdrawalRepository
	depositDatabase    deposit.DepositRepository
	transferDatabase   transfer.TransferRepository
	autodebitDatabase  automaticdebit.AutoDebitRepository
}

func NewAccountService(accountDB account.AccountRepository,
	userServe UserService,
	clientServe ClientService,
	withdrawalDB withdrawal.WithdrawalRepository,
	depositDB deposit.DepositRepository,
	transferDB transfer.TransferRepository,
	autodebitDB automaticdebit.AutoDebitRepository,
) AccountService {
	return accountServiceImpl{
		accountDatabase:    accountDB,
		userService:        userServe,
		clientService:      clientServe,
		withdrawalDatabase: withdrawalDB,
		depositDatabase:    depositDB,
		transferDatabase:   transferDB,
		autodebitDatabase:  autodebitDB,
	}
}

func (service accountServiceImpl) Create(ctx context.Context, accountRequest *account.Account) (*account.Account, *model.Erro) {
	if accountRequest.User_id == "" || accountRequest.Client_id == "" {
		log.Warn().Msg("Missing credentials on creating account")
		return nil, ErrorMissingCredentials
	}
	// verify if client and user exists, PERMISSION MUST BE of user
	if _, err := service.userService.Get(ctx, &accountRequest.User_id); err == model.IDnotFound || err != nil {
		return nil, err
	}

	if _, err := service.clientService.Get(ctx, &accountRequest.Client_id); err == model.IDnotFound || err != nil {
		return nil, err
	}

	accountResponse, err := service.accountDatabase.Create(ctx, accountRequest)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Account created")
	return accountResponse, nil
}

func (service accountServiceImpl) Get(ctx context.Context, accountID *string) (*account.Account, *model.Erro) {
	accountResponse, err := service.accountDatabase.Get(ctx, accountID)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Account returned: " + *accountID)
	return accountResponse, nil
}

func (service accountServiceImpl) Update(ctx context.Context, accountRequest *account.Account) (*account.Account, *model.Erro) {
	accountResponse, err := service.Get(ctx, &accountRequest.Account_id)
	if err != nil {
		return nil, err
	}

	// verifica valores que foram passados ou não
	if accountRequest.Account_id == "" {
		log.Warn().Msg("No account with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("account id invalid"), HttpCode: http.StatusBadRequest}
	}
	if accountRequest.Agency_id != 0 {
		accountResponse.Agency_id = accountRequest.Agency_id
	}
	if accountRequest.Client_id != "" {
		accountResponse.Client_id = accountRequest.Client_id
	}
	if accountRequest.User_id != "" {
		accountResponse.User_id = accountRequest.User_id
	}
	if accountRequest.Balance != accountResponse.Balance {
		accountResponse.Balance = accountRequest.Balance
	}
	if accountRequest.Status != "" {
		if !(accountRequest.Status).IsValid() {
			return nil, model.InvalidStatus
		}
		accountResponse.Status = accountRequest.Status
	}
	if err := service.accountDatabase.Update(ctx, accountResponse); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (account): " + accountRequest.Account_id)

	return service.Get(ctx, &accountRequest.Account_id)
}

func (service accountServiceImpl) GetAll(ctx context.Context) (*[]account.Account, *model.Erro) {
	accountSlice, err := service.accountDatabase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return accountSlice, nil
}

func (service accountServiceImpl) Delete(ctx context.Context, accountID *string) *model.Erro {
	if err := service.accountDatabase.Delete(ctx, accountID); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + *accountID)
	return nil
}

/*
func (service accountServiceImpl) Status(accountID *string, status bool) *model.Erro {
	account, err := service.Get(accountID)
	if err != nil {
		return err
	}
	account.Status = status
	if err := service.accountDatabase.Update(account); err != nil {
		return err
	}
	log.Info().Msg("Account status changed: " + *accountID + " to " + strconv.FormatBool(status))
	return nil
}
*/

func (service accountServiceImpl) Report(ctx context.Context, accountID *string) (*account.AccountReport, *model.Erro) {
	accountInfo, err := service.Get(ctx, accountID)
	if err != nil {
		return nil, err
	}
	transfers, err := service.transferDatabase.GetFilteredByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	deposits, err := service.depositDatabase.GetFilteredByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	withdrawals, err := service.withdrawalDatabase.GetFilteredByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	automaticDebits, err := service.autodebitDatabase.GetFilteredByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	accountReport := account.AccountReport{
		Account_id:       accountInfo.Account_id,
		Client_id:        accountInfo.Client_id,
		Agency_id:        accountInfo.Agency_id,
		Balance:          accountInfo.Balance,
		Register_date:    accountInfo.Register_date,
		Status:           accountInfo.Status,
		Transfers:        transfers,
		Deposits:         deposits,
		Withdrawals:      withdrawals,
		Automatic_Debits: automaticDebits,
		Report_Date:      time.Now().Format(timeLayout),
	}
	log.Info().Msg("Report generated for account: " + *accountID)
	return &accountReport, nil
}
