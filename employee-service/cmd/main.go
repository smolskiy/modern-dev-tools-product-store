// cmd/main.go
package main

import (
	"log"
	"net/http"

	handler "catalog-service/internal/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/employees", handler.GetEmployees).Methods("GET")
	r.HandleFunc("/employees", handler.CreateEmployee).Methods("POST")
	// ... остальные роуты
	log.Println("Employee service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
