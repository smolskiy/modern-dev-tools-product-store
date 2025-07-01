package handler

import (
	model "catalog-service/internal/models"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

var (
	users = map[string]model.User{
		"admin": {ID: 1, Username: "admin", Password: "admin", Role: "admin"},
		"hr":    {ID: 2, Username: "hr", Password: "hr", Role: "hr"},
	}
	mu     sync.Mutex
	jwtKey = []byte("your-secret-key")
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, ok := users[creds.Username]
	mu.Unlock()
	if !ok || user.Password != creds.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[creds.Username]; exists {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}
	id := len(users) + 1
	users[creds.Username] = model.User{ID: id, Username: creds.Username, Password: creds.Password, Role: creds.Role}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}
