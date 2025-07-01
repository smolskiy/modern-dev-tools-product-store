package main

import (
	"log"
	"net/http"

	handler "catalog-service/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// Должности
	r.HandleFunc("/positions", handler.GetPositions).Methods("GET")
	r.HandleFunc("/positions", handler.CreatePosition).Methods("POST")
	r.HandleFunc("/positions/{id}", handler.UpdatePosition).Methods("PUT")
	r.HandleFunc("/positions/{id}", handler.DeletePosition).Methods("DELETE")
	// Добавить аналогично квалификации и типы отпусков...
	log.Println("Directory service started on :8084")
	log.Fatal(http.ListenAndServe(":8084", r))
}
