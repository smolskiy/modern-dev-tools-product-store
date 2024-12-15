package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

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
}

// Product представляет структуру товара
type Product struct {
	ID    int     `xml:"id"`
	Name  string  `xml:"name"`
	Price float64 `xml:"price"`
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	SoapEnv string   `xml:"xmlns:soapenv,attr"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	XMLName   xml.Name     `xml:"Body"`
	ProductOp ProductOpReq `xml:"ProductOperation"`
}

type ProductOpReq struct {
	XMLName xml.Name `xml:"ProductOperation"`
	Action  string   `xml:"action"`
	Product Product  `xml:"product"`
}

// SOAPResponse представляет ответ SOAP
type SOAPResponse struct {
	XMLName xml.Name `xml:"ProductOperationResponse"`
	Message string   `xml:"message"`
	Product *Product `xml:"product,omitempty"`
}

func soapHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received SOAP request:")
	body, _ := io.ReadAll(r.Body)
	log.Println(string(body))

	var envelope SOAPEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		log.Printf("Error decoding SOAP request: %v", err)
		http.Error(w, fmt.Sprintf("Invalid SOAP request: %v", err), http.StatusBadRequest)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		log.Printf("Error decoding SOAP request: %v", err)
		http.Error(w, fmt.Sprintf("Invalid SOAP request: %v", err), http.StatusBadRequest)
		return
	}

	action := envelope.Body.ProductOp.Action
	product := envelope.Body.ProductOp.Product

	var response SOAPResponse
	switch action {
	case "add":
		err := db.QueryRow("INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
			product.Name, product.Price).Scan(&product.ID)
		if err != nil {
			http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
			return
		}
		response.Message = "Product added successfully"
		response.Product = &product

	case "get":
		row := db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", product.ID)
		err := row.Scan(&product.ID, &product.Name, &product.Price)
		if err == sql.ErrNoRows {
			response.Message = "Product not found"
		} else if err != nil {
			http.Error(w, "Ошибка чтения из базы данных", http.StatusInternalServerError)
			return
		} else {
			response.Message = "Product retrieved successfully"
			response.Product = &product
		}

	case "getAll":
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

		response.Message = "Products retrieved successfully"
		responseXML, _ := xml.Marshal(products)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(responseXML)

	default:
		response.Message = "Invalid action"
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(SOAPEnvelope{
		SoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SOAPBody{
			ProductOp: ProductOpReq{
				Action:  "response",
				Product: Product{},
			},
		},
	})
}

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установка заголовков CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Предзапрос OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting SOAP service...")

	initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/soap", soapHandler).Methods("POST")

	log.Println("SOAP server started on :8080")
	if err := http.ListenAndServe(":8080", withCORS(router)); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
