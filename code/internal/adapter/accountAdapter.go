package adapter

import (
	"net/http"
	"strconv"

	adapter "BankingAPI/code/internal/adapter/infrastructure"
	"BankingAPI/code/internal/application"
	"BankingAPI/code/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func AccountPostAdapter(c *echo.Context) (int, interface{}) {
	var accountInfo domain.Account

	userID, err := strconv.ParseUint((*c).FormValue("UserID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	clientID, err := strconv.ParseUint((*c).FormValue("ClientID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	agencyID, err := strconv.ParseUint((*c).FormValue("AgencyID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	accountInfo.SetAccountClientId(uint32(clientID))
	accountInfo.SetAccountUserId(uint32(userID))
	accountInfo.SetAgencyID(uint32(agencyID))
	accountInfo.SetAccountPassword((*c).FormValue("Password"))
	accountInfo.SetAccountBalance(0.0)
	accountInfo.SetStatus(true)

	err = adapter.CreateAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, accountInfo
}

func AccountGetAdapter(c *echo.Context) (int, interface{}) {
	var accountInfo domain.Account

	accountID, err := strconv.ParseUint((*c).FormValue("AccountID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	accountInfo.SetAccountId(uint32(accountID))

	err = adapter.GetAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, accountInfo
}

// PRECISA CHAMAR APPLICATION
func AccountPutAdapter(c *echo.Context) (int, interface{}) {
	var accountInfo domain.Account

	accountID, err := strconv.Atoi((*c).Param("AccountID"))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	userID, err := strconv.ParseUint((*c).FormValue("UserID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	clientID, err := strconv.ParseUint((*c).FormValue("ClientID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	agencyID, err := strconv.ParseUint((*c).FormValue("AgencyID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	accountInfo.SetAccountId(uint32(accountID))
	accountInfo.SetAccountClientId(uint32(clientID))
	accountInfo.SetAccountUserId(uint32(userID))
	accountInfo.SetAgencyID(uint32(agencyID))
	accountInfo.SetAccountPassword((*c).FormValue("Password"))

	err = adapter.UpdateAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, accountInfo
}

func AccountDeleteAdapter(c *echo.Context) (int, interface{}) {
	accountID, err := strconv.ParseUint((*c).Param("AccountID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	err = adapter.DeleteAccountDB(uint32(accountID))
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, domain.Account{AccountId: uint32(accountID)}
}

func AccountGetByOrderFilterAdapter(c *echo.Context) (int, interface{}) {
	queryParams := (*c).QueryParams()
	if !(queryParams.Has("filter") && (queryParams.Has("order"))) {
		err := "Missing parameters values: filter:" + fromBoolToString(queryParams.Has("filter")) + " order:" + fromBoolToString(queryParams.Has("order"))
		log.Warn().Msg(err)
		return http.StatusBadRequest, domain.ErrorResponse{Message: err}
	}
	// filter by bank ID order by accountID
	listOfAccounts := make([]domain.Account, 0, 0)

	err := adapter.SearchAccountByFilterOrder(queryParams.Get("filter"), queryParams.Get("order"), &listOfAccounts)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	return http.StatusOK, listOfAccounts
}

func fromBoolToString(value bool) string {
	return strconv.FormatBool(value)
}

func AccountPutDepositAdapter(c *echo.Context) (int, interface{}) {
	var accountInfo domain.Account
	accountID, err := strconv.ParseUint((*c).Param("AccountID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}
	accountInfo.SetAccountId(uint32(accountID))

	balanceUpdate, err := strconv.ParseFloat((*c).FormValue("BalanceUpdate"), 64)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	err = adapter.GetAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	err = application.Deposit(&accountInfo, balanceUpdate)
	if err != nil {
		log.Warn().Msg(err.Error())
		return http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()}
	}

	err = adapter.UpdateAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, domain.Account{AccountId: accountInfo.GetAccountId(), AccountBalance: accountInfo.GetBalance()}
}

func AccountPutWithdrawalAdapter(c *echo.Context) (int, interface{}) {
	var accountInfo domain.Account
	accountID, err := strconv.ParseUint((*c).Param("AccountID"), 0, 32)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	accountInfo.SetAccountId(uint32(accountID))

	balanceUpdate, err := strconv.ParseFloat((*c).FormValue("BalanceUpdate"), 64)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	err = adapter.GetAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	err = application.Withdraw(&accountInfo, balanceUpdate)
	if err != nil {
		log.Warn().Msg(err.Error())
		return http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()}
	}

	err = adapter.UpdateAccountDB(&accountInfo)
	if err != nil {
		log.Error().Msg(err.Error())
		return http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()}
	}

	return http.StatusOK, domain.Account{AccountId: accountInfo.GetAccountId(), AccountBalance: accountInfo.GetBalance()}
}
