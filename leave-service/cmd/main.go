package main

import (
	"log"
	"net/http"

	handler "catalog-service/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/leaves", handler.GetLeaves).Methods("GET")
	r.HandleFunc("/leaves", handler.CreateLeave).Methods("POST")
	r.HandleFunc("/leaves/{id}/status", handler.UpdateLeaveStatus).Methods("PUT")
	log.Println("Leave service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
