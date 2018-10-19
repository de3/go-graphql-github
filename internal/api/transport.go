package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHandler(apiService Service) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", apiService.command).Methods("POST")

	return r
}
