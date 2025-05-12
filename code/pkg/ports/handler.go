package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type User struct {
	Id           int32  `json:"id" xml:"id"`
	Name         string `json:"name" xml:"name"`
	RegisterDate string `json:"register_date" xml:"register_date"`
}

func Server() {
	http.HandleFunc("/user/{userID}", userResponse)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}

func userResponse(w http.ResponseWriter, r *http.Request) {
	user := User{
		Id:           (int32(1)),
		Name:         "julio",
		RegisterDate: "2016-08-29T09:12:33.001Z",
	}
	erro := errors.New("rolou um problema")
	log.Info().Any("user", user).AnErr("erro", erro).Msg("Request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Debug
// InfoLevel
// WarnLevel
// ErrorLevel
// FatalLevel
// PanicLevel
