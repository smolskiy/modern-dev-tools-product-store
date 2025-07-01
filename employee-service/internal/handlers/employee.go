// internal/handler/employee.go
package handler

import (
	model "catalog-service/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
)

var (
	employees = make(map[int]model.Employee)
	nextID    = 1
	mu        sync.Mutex
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	list := make([]model.Employee, 0, len(employees))
	for _, emp := range employees {
		list = append(list, emp)
	}
	json.NewEncoder(w).Encode(list)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	emp.ID = nextID
	nextID++
	employees[emp.ID] = emp
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

// ... другие CRUD-операции (см. выше)
