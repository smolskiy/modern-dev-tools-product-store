package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	model "catalog-service/internal/models"
	"github.com/gorilla/mux"
)

var (
	departments = make(map[int]model.Department)
	nextID      = 1
	mu          sync.Mutex
)

func GetDepartments(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	list := make([]model.Department, 0, len(departments))
	for _, d := range departments {
		list = append(list, d)
	}
	json.NewEncoder(w).Encode(list)
}

func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var dep model.Department
	if err := json.NewDecoder(r.Body).Decode(&dep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	dep.ID = nextID
	nextID++
	departments[dep.ID] = dep
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dep)
}

func GetDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	dep, ok := departments[id]
	mu.Unlock()
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(dep)
}

func UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var dep model.Department
	if err := json.NewDecoder(r.Body).Decode(&dep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	dep.ID = id
	departments[id] = dep
	mu.Unlock()
	json.NewEncoder(w).Encode(dep)
}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	delete(departments, id)
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
