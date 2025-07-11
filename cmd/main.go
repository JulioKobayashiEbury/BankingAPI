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
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "banking")

	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	gcproject := os.Getenv("GOOGLE_CLOUD_PROJECT")

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
