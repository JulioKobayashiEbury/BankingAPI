package ports

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id  int32 `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	RegisterDate string `json:"register_date" xml:"register_date"`
}

func Server() {
	http.HandleFunc("/user/{userID}", userResponse)
	http.ListenAndServe(":8080", nil)
}

func userResponse(w http.ResponseWriter, r *http.Request) {
	user := User{
		Id:  (int32(1)),
		Name: "julio",
		RegisterDate: "2016-08-29T09:12:33.001Z",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
