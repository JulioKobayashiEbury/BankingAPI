package controller

/* zerolog.SetGlobalLevel(zerolog.InfoLevel)
log.Info().Msg("Method not allowed")
w.WriteHeader(http.StatusMethodNotAllowed)
response := map[string]string{"error": "Method not allowed"}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
return
*/

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	documentLenghtIdeal = 14
	maxNameLenght       = 30
)

var DatabaseClient *firestore.Client

func Server() {
	server := echo.New()
	AddAccountEndPoints(server)
	AddClientsEndPoints(server)
	AddUsersEndPoints(server)

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if err := server.Start("localhost:25565"); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg("Server started on port 25565")
}
