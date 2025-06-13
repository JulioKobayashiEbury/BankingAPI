package service

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	model "BankingAPI/internal/model"
	"BankingAPI/internal/model/account"

	"github.com/rs/zerolog/log"
)

type accountServiceImpl struct {
	userService     UserService
	clientService   ClientService
	accountDatabase model.RepositoryInterface

	withdrawalDatabase model.RepositoryInterface
	depositDatabase    model.RepositoryInterface
	transferDatabase   model.RepositoryInterface
	autodebitDatabase  model.RepositoryInterface
}

func NewAccountService(accountDB model.RepositoryInterface,
	userServe UserService,
	clientServe ClientService,
	withdrawalDB model.RepositoryInterface,
	depositDB model.RepositoryInterface,
	transferDB model.RepositoryInterface,
	autodebitDB model.RepositoryInterface) AccountService {
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

func (service accountServiceImpl) Create(accountRequest *account.Account) (*account.Account, *model.Erro) {
	if accountRequest.User_id == "" || accountRequest.Client_id == "" {
		log.Warn().Msg("Missing credentials on creating account")
		return nil, ErrorMissingCredentials
	}
	// verify if client and user exists, PERMISSION MUST BE of user
	if _, err := service.userService.Get(&accountRequest.User_id); err == model.IDnotFound || err != nil {
		return nil, err
	}

	if _, err := service.clientService.Get(&accountRequest.Client_id); err == model.IDnotFound || err != nil {
		return nil, err
	}

	obj, err := service.accountDatabase.Create(accountRequest)
	if err != nil {
		return nil, err
	}
	accountResponse, ok := obj.(*account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}

	log.Info().Msg("Account created")
	return accountResponse, nil
}

func (service accountServiceImpl) Get(accountID *string) (*account.Account, *model.Erro) {
	obj, err := service.accountDatabase.Get(accountID)
	if err != nil {
		return nil, err
	}
	accountResponse, ok := obj.(*account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	log.Info().Msg("Account returned: " + *accountID)
	return accountResponse, nil
}

func (service accountServiceImpl) Update(accountRequest *account.Account) (*account.Account, *model.Erro) {
	accountResponse, err := service.Get(&accountRequest.Account_id)
	if err != nil {
		return nil, err
	}

	// verifica valores que foram passados ou não
	if accountRequest.Account_id == "" {
		log.Warn().Msg("No account with id: 0 allowed")
		return nil, &model.Erro{Err: errors.New("Account id invalid"), HttpCode: http.StatusBadRequest}
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
	// monta struct de update

	if err := service.accountDatabase.Update(accountResponse); err != nil {
		return nil, err
	}

	log.Info().Msg("Update was succesful (account): " + accountRequest.Account_id)

	return service.Get(&accountRequest.Account_id)
}

func (service accountServiceImpl) GetAll() (*[]account.Account, *model.Erro) {
	obj, err := service.accountDatabase.GetAll()
	if err != nil {
		return nil, err
	}
	accountSlice, ok := obj.(*[]account.Account)
	if !ok {
		return nil, model.DataTypeWrong
	}
	return accountSlice, nil
}

func (service accountServiceImpl) Delete(accountID *string) *model.Erro {
	if err := service.accountDatabase.Delete(accountID); err != nil {
		return err
	}
	log.Info().Msg("Account deleted: " + *accountID)
	return nil
}

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

func (service accountServiceImpl) Report(accountID *string) (*account.AccountReport, *model.Erro) {
	accountInfo, err := service.Get(accountID)
	if err != nil {
		return nil, err
	}
	filters := []string{"account_id,==," + *accountID}

	transfers, err := service.transferDatabase.GetFiltered(&filters)
	if err != nil {
		return nil, err
	}

	deposits, err := service.depositDatabase.GetFiltered(&filters)
	if err != nil {
		return nil, err
	}
	withdrawals, err := service.withdrawalDatabase.GetFiltered(&filters)
	if err != nil {
		return nil, err
	}
	automaticDebits, err := service.autodebitDatabase.GetFiltered(&filters)
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
