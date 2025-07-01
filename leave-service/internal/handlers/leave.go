package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	model "catalog-service/internal/models"
	"github.com/gorilla/mux"
)

var (
	leaves = make(map[int]model.Leave)
	nextID = 1
	mu     sync.Mutex
)

func GetLeaves(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	list := make([]model.Leave, 0, len(leaves))
	for _, l := range leaves {
		list = append(list, l)
	}
	json.NewEncoder(w).Encode(list)
}

func CreateLeave(w http.ResponseWriter, r *http.Request) {
	var leave model.Leave
	if err := json.NewDecoder(r.Body).Decode(&leave); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	leave.ID = nextID
	leave.Status = "заявка"
	leave.CreatedAt = time.Now()
	nextID++
	leaves[leave.ID] = leave
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(leave)
}

func UpdateLeaveStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	leave, ok := leaves[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	leave.Status = payload.Status
	leaves[id] = leave
	mu.Unlock()
	json.NewEncoder(w).Encode(leave)
}
