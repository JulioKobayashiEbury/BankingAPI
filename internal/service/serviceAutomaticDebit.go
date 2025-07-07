package service

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	model "BankingAPI/internal/model"
	automaticdebit "BankingAPI/internal/model/automaticDebit"
	"BankingAPI/internal/model/withdrawal"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

const (
	//"2006-01-02T15:04:05+07:00"
	timeLayout = time.RFC3339
)

type serviceAutoDebitImpl struct {
	autoDebitFirestore automaticdebit.AutoDebitRepository
	withdrawalService  WithdrawalService
}

func NewAutoDebit(autodebitDB automaticdebit.AutoDebitRepository, withdrawal WithdrawalService) AutomaticDebitService {
	return serviceAutoDebitImpl{
		autoDebitFirestore: autodebitDB,
		withdrawalService:  withdrawal,
	}
}

func (service serviceAutoDebitImpl) Create(ctx context.Context, autodebitRequest *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *echo.HTTPError) {
	automaticDebit, err := service.autoDebitFirestore.Create(ctx, autodebitRequest)
	if err != nil {
		return nil, err
	}
	return automaticDebit, nil
}

func (service serviceAutoDebitImpl) Delete(ctx context.Context, id *string) *echo.HTTPError {
	if err := service.autoDebitFirestore.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (service serviceAutoDebitImpl) Get(ctx context.Context, id *string) (*automaticdebit.AutomaticDebit, *echo.HTTPError) {
	autodebit, err := service.autoDebitFirestore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return autodebit, nil
}

func (service serviceAutoDebitImpl) GetAll(ctx context.Context) (*[]automaticdebit.AutomaticDebit, *echo.HTTPError) {
	autodebits, err := service.autoDebitFirestore.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return autodebits, nil
}

func (service serviceAutoDebitImpl) ProcessNewAutomaticDebit(ctx context.Context, autoDebit *automaticdebit.AutomaticDebit) (*automaticdebit.AutomaticDebit, *echo.HTTPError) {
	if !isValidDate(autoDebit.Expiration_date) {
		log.Warn().Msg("Invalid date format")
		return nil, &echo.HTTPError{Internal: errors.New("invalid date format"), Code: http.StatusBadRequest, Message: "invalid date format"}
	}
	autodebitResponse, err := service.autoDebitFirestore.Create(ctx, autoDebit)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Automatic debit created: " + autodebitResponse.Debit_id)
	return autodebitResponse, nil
}

func isValidDate(date string) bool {
	_, err := time.Parse(model.TimeLayout, date)
	return err == nil
}

func (service serviceAutoDebitImpl) CheckAutomaticDebits() {
	log.Info().Msg("Checking for auto debits...")
	ctx := context.Background()
	defer ctx.Done()
	autoDebitList, err := service.GetAll(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	for _, autoDebit := range *autoDebitList {
		logger := log.With().Fields(autoDebit).Logger()
		expirationDate, err := time.Parse(timeLayout, autoDebit.Expiration_date)
		if err != nil {
			logger.Error().Msg(err.Error())
			return
		}
		if expirationDate.Unix() < time.Now().Unix() {
			if err := service.Delete(ctx, &autoDebit.Debit_id); err != nil {
				logger.Error().Msg("Failed to delete expired automatic debit")
				return
			}
			logger.Warn().Msg("Debit is expired, deleted this automatic debit")
			continue
		}
		if autoDebit.Debit_day == uint16(time.Now().Day()) {
			_, err := service.withdrawalService.ProcessWithdrawal(ctx, &withdrawal.Withdrawal{
				Account_id: autoDebit.Account_id,
				User_id:    autoDebit.User_id,
				Agency_id:  autoDebit.Agency_id,
				Withdrawal: autoDebit.Value,
			})
			if err != nil {
				logger.Error().Msg(err.Error())
				return
			}
			logger.Info().Msg("Auto debit is logged as Withdrawal together with Account")
		}
	}
}

func (service serviceAutoDebitImpl) Scheduled() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return
	}

	job, err := scheduler.NewJob(
		/*
			gocron.DailyJob(1, gocron.NewAtTimes(
				gocron.NewAtTime(10, 00, 00)
			))
		*/
		gocron.CronJob(
			"*/2 * * * *",
			false,
		),
		gocron.NewTask(
			service.CheckAutomaticDebits,
		),
		gocron.WithName("Checking Automatic Debits"),
	)
	if err != nil {
		return
	}
	scheduler.Start()
	log.Info().Msg("Job ID for automatic debit checking: " + job.ID().String())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info().Msg("Interruption signal received, terminating gracefully...")
		scheduler.Shutdown()
		os.Exit(0)
	}()
}
