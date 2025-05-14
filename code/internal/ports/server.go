package ports

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

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Server() {
	server := echo.New()
	// User
	server.POST("/user", UserPostHandler)
	server.GET("/user/:UserID", UserGetHandler)
	server.DELETE("/user/:UserID", UserDeleteHandler)
	server.PUT("/user/:UserID", UserPutHandler)

	// Client

	// Account

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Server started on port 8080")
	server.Start("localhost:8080")
}

// Debug
// InfoLevel
// WarnLevel
// ErrorLevel
// FatalLevel
// PanicLevel
