package ports

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func Server() {
	http.HandleFunc("/user", userResponse)
	http.ListenAndServe(":8080", nil)
}

func userResponse(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name:  "Julio",
		Email: "julio@gmail.com",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
