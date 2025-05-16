package initializers

import (
	"context"

	"github.com/rs/zerolog/log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	FireBaseApp *firebase.App
	Ctx         context.Context
)

// initalaize DB conn
func InitializeDB() (*firebase.App, error) {
	Ctx := context.Background()

	// Initialize Firebase Admin SDK
	// Replace "path/to/your/serviceAccountKey.json" with the actual path
	sa := option.WithCredentialsFile("path/to/your/serviceAccountKey.json")
	FireBaseApp, err := firebase.NewApp(Ctx, nil, sa)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return FireBaseApp, nil
}
