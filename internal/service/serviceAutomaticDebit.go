package service

import (
	"net/http"
	"time"

	repository "BankingAPI/internal/model/repository"
	model "BankingAPI/internal/model/types"

	"github.com/rs/zerolog/log"
)

const timeLayout = "2006-01-02 15:04:05.999999999 -0700"

func ProcessNewAutomaticDebit(autoDebit *model.AutomaticDebitRequest) *model.Erro {
	// map new debit
	times, err := time.Parse(timeLayout, (*autoDebit).Debit_date)
	if err != nil {
		return &model.Erro{Err: err, HttpCode: http.StatusInternalServerError}
	}
	debitMap := map[string]interface{}{
		"account_id":    autoDebit.Account_id,
		"client_id":     autoDebit.Client_id,
		"agency_id":     autoDebit.Agency_id,
		"value":         autoDebit.Value,
		"debit_date":    times.UnixMicro(),
		"register_date": time.Now().UnixMicro(),
	}
	// insert in db
	var debitId string
	if err := repository.CreateObject(&debitMap, repository.AutoDebit, &debitId); err != nil {
		return err
	}
	log.Info().Msg("Automatic debit created: " + debitId)
	return nil
}

func scheduleNewDebit() {
}

func CheckAutomaticDebits() {
}
