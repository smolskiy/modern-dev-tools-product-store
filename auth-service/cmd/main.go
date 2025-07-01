package main

import (
	"log"
	"net/http"

	handler "catalog-service/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	log.Println("Auth service started on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
