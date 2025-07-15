package main

import (
	"context"
	"io"
	"os"
	"time"

	"BankingAPI/internal/controller"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/service"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logFileName := os.Getenv("LOG_FILE")
	if logFileName == "" {
		log.Fatal().Msg("LOG_FILE environment variable not set")
		return
	}
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		log.Fatal().Msg("FIRESTORE_EMULATOR_HOST environment variable not set")
		return
	}

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	gcproject := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if gcproject == "" {
		log.Fatal().Msg("GOOGLE_CLOUD_PROJECT environment variable not set")
		return
	}

	ctx := context.Background()
	defer ctx.Done()

	client, err := firestore.NewClient(ctx, gcproject)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	repositories := controller.InstantiateRepo(client)
	gateways := controller.InstantiateGateways()

	services := controller.InstantiateServices(repositories, gateways)

	if err := createAdminUser(services.UserService); err != nil {
		log.Panic().Msg("Not able to create admin user!")
		return
	}

	controller.Server(services)
}

func createAdminUser(userService service.UserService) error {
	ctx := context.Background()
	defer ctx.Done()

	allUsers, err := userService.GetAll(ctx)
	if err != nil {
		return err.Internal
	}

	for _, user := range *allUsers {
		if user.Name == "admin" {
			log.Info().Msg("Admin user already exists with ID: " + user.User_id)
			return nil
		}
	}

	userResponse, err := userService.Create(ctx, &user.User{
		Name:     "admin",
		Document: "00000000000000",
		Password: "admin",
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return err.Internal
	}
	log.Info().Msg("Admin user created with ID: " + userResponse.User_id)
	return nil
}
