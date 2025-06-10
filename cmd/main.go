package main

import (
	"context"
	"os"

	"BankingAPI/internal/controller"
	"BankingAPI/internal/model/user"
	"BankingAPI/internal/service"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func main() {
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

	repositories := controller.InstantiateRepo(client)
	services := controller.InstantiateServices(repositories)

	if err := createAdminUser(services.UserService); err != nil {
		log.Panic().Msg("Not able to create admin user!")
		return
	}

	controller.Server(services)
}

func createAdminUser(userService service.UserService) error {
	allUsers, err := userService.GetAll()
	if err != nil {
		return err.Err
	}

	for _, user := range *allUsers {
		if user.Name == "admin" {
			log.Info().Msg("Admin user already exists with ID: " + user.User_id)
			return nil
		}
	}

	userResponse, err := userService.Create(&user.User{
		Name:     "admin",
		Document: "00000000000000",
		Password: "admin",
	})
	if err != nil {
		log.Error().Msg(err.Err.Error())
		return err.Err
	}
	log.Info().Msg("Admin user created with ID: " + userResponse.User_id)
	return nil
}
