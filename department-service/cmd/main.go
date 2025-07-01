package main

import (
	"log"
	"net/http"

	handler "catalog-service/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/departments", handler.GetDepartments).Methods("GET")
	r.HandleFunc("/departments", handler.CreateDepartment).Methods("POST")
	r.HandleFunc("/departments/{id}", handler.GetDepartment).Methods("GET")
	r.HandleFunc("/departments/{id}", handler.UpdateDepartment).Methods("PUT")
	r.HandleFunc("/departments/{id}", handler.DeleteDepartment).Methods("DELETE")
	log.Println("Department service started on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
