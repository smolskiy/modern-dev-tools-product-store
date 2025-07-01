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
	positions = make(map[int]model.Position)
	posID     = 1
	mu        sync.Mutex
)

func GetPositions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	list := make([]model.Position, 0, len(positions))
	for _, p := range positions {
		list = append(list, p)
	}
	json.NewEncoder(w).Encode(list)
}

func CreatePosition(w http.ResponseWriter, r *http.Request) {
	var p model.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	p.ID = posID
	posID++
	positions[p.ID] = p
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func UpdatePosition(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var p model.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	p.ID = id
	positions[id] = p
	mu.Unlock()
	json.NewEncoder(w).Encode(p)
}

func DeletePosition(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	delete(positions, id)
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
