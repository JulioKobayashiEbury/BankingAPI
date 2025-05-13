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
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Server() {
	// POST HANDLERS
	http.HandleFunc("/user", userPostHandler)
	http.HandleFunc("/account", accountPostHandler)
	http.HandleFunc("/client", clientPostHandler)
	// OTHER METHODS
	http.HandleFunc("/user/{UserID}", userHandler)
	http.HandleFunc("/account/{AccountID}", accountHandler)
	http.HandleFunc("/client/{ClientID}", clientHandler)

	// SPECIFIC HANDLERS
	http.HandleFunc("user/{UserID}/block", userBlockHandler)
	http.HandleFunc("user/{UserID}/unblock", userUnBlockHandler)
	http.HandleFunc("clients/{ClientID}/block", clientBlockHandler)
	http.HandleFunc("clients/{ClientID}/unblock", clientUnBlockHandler)
	http.HandleFunc("accounts/{AccountID}/block", accountBlockHandler)
	http.HandleFunc("accounts/{AccountID}/unblock", accountUnBlockHandler)

	// Automatic Debit
	http.HandleFunc("accounts/{AccountID}/automaticDebit", accountDebitHandler)

	// Balance
	http.HandleFunc("accounts/{AccountID}/balance/deposit", accountBalanceDepositHandler)
	http.HandleFunc("accounts/{AccountID}/balance/withdraw", accountBalanceWithdrawHandler)
	http.HandleFunc("accounts/", accountGetFilteredOrderedHandler)
	http.HandleFunc("accounts/{AccountID}/newTransfer", accountNewTranferHandler)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

func accountNewTranferHandler(w http.ResponseWriter, r *http.Request) {
}

func accountBalanceDepositHandler(w http.ResponseWriter, r *http.Request) {
}

func accountBalanceWithdrawHandler(w http.ResponseWriter, r *http.Request) {
}

func accountGetFilteredOrderedHandler(w http.ResponseWriter, r *http.Request) {
}

func accountDebitHandler(w http.ResponseWriter, r *http.Request) {
}

func userBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func userUnBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func clientBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func clientUnBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func accountBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func accountUnBlockHandler(w http.ResponseWriter, r *http.Request) {
}

func userHandler(w http.ResponseWriter, r *http.Request) {
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
}

// POST HANDLERS
// userPostHandler handles POST requests for creating a new user
// accountPostHandler handles POST requests for creating a new account
// clientPostHandler handles POST requests for creating a new client
func userPostHandler(w http.ResponseWriter, r *http.Request) {
}

func accountPostHandler(w http.ResponseWriter, r *http.Request) {
}

func clientPostHandler(w http.ResponseWriter, r *http.Request) {
}

// Debug
// InfoLevel
// WarnLevel
// ErrorLevel
// FatalLevel
// PanicLevel
