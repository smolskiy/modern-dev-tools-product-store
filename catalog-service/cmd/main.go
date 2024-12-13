package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// Product представляет структуру товара
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var db *sql.DB

func initDB() {
	var err error
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Создание таблицы, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price NUMERIC(10, 2) NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		log.Fatalf("Ошибка проверки данных в таблице: %v", err)
	}

	// Если таблица пуста, вставляем одну запись
	if count == 0 {
		log.Println("Таблица пустая. Добавляем начальную запись.")
		_, err := db.Exec("INSERT INTO products (name, price) VALUES ($1, $2)",
			"Product 1", 100.0)
		if err != nil {
			log.Fatalf("Ошибка вставки начальной записи: %v", err)
		}
	}
}

// Обработчик для получения всех товаров
func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, "Ошибка чтения из базы данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			http.Error(w, "Ошибка сканирования данных", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Обработчик для добавления нового товара
func addProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
		product.Name, product.Price).Scan(&product.ID)
	if err != nil {
		http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Обработчик для обновления товара
func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлекаем параметры из URL
	idStr := vars["id"] // Получаем ID из URL

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный формат ID", http.StatusBadRequest)
		return
	}
	// Декодируем тело запроса в структуру Product
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Выполняем запрос на обновление
	_, err = db.Exec("UPDATE products SET name = $1, price = $2 WHERE id = $3",
		product.Name, product.Price, id)
	if err != nil {
		http.Error(w, "Ошибка обновления данных в базе", http.StatusInternalServerError)
		return
	}

	// Возвращаем обновленные данные
	product.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Обработчик для удаления товара
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлекаем параметры из URL
	id := vars["id"]    // Получаем ID из URL

	// Выполняем запрос на удаление
	_, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка удаления данных из базы", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Отправляем статус 204 (без контента)
}

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting catalog service...")

	initDB()
	defer db.Close()

	router := mux.NewRouter()

	// Определяем маршруты
	router.HandleFunc("/products", getProducts).Methods("GET")
	router.HandleFunc("/products", addProduct).Methods("POST")
	router.HandleFunc("/products/{id:[0-9]+}", updateProduct).Methods("PUT")    // Метод PUT для обновления
	router.HandleFunc("/products/{id:[0-9]+}", deleteProduct).Methods("DELETE") // Метод DELETE для удаления

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", withCORS(router)); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
