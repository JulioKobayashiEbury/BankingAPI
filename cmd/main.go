package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"BankingAPI/internal/controller"

	"cloud.google.com/go/firestore"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

func init() {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "0.0.0.0:8080")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "banking")

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()
	defer ctx.Done()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	controller.DatabaseClient = client

	controller.InstantiateRepo()
	controller.InstantiateServices()
}

func main() {
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
			controller.Services.AutomaticdebitService.CheckAutomaticDebits,
		),
		gocron.WithName("Checking Automatic Debits"),
	)
	if err != nil {
		return
	}
	scheduler.Start()
	log.Info().Msg(job.ID().String())
	fmt.Print("Scheduler running...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info().Msg("Interruption signal received, terminating gracefully...")
		scheduler.Shutdown()
		os.Exit(0)
	}()
	controller.Server()
}
