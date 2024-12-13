package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func callSOAPService(action string, payload string) (string, error) {
	soapURL := "http://host.docker.internal:8082/soap"

	// SOAP-запрос
	soapRequest := fmt.Sprintf(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
	   <soapenv:Body>
		  <ProductOperation>
			 <action>%s</action>
			 %s
		  </ProductOperation>
	   </soapenv:Body>
	</soapenv:Envelope>`, action, payload)

	// Отправка запроса
	resp, err := http.Post(soapURL, "text/xml", bytes.NewBufferString(soapRequest))
	if err != nil {
		return "", fmt.Errorf("Ошибка отправки SOAP-запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения ответа SOAP-сервиса: %v", err)
	}

	return string(responseData), nil
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	// Формируем тело SOAP-запроса
	payload := fmt.Sprintf("<product><id>%s</id></product>", id)

	// Отправляем SOAP-запрос
	response, err := callSOAPService("get", payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ клиенту
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(response))
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем данные из тела запроса
	name := r.URL.Query().Get("name")
	price := r.URL.Query().Get("price")
	if name == "" || price == "" {
		http.Error(w, "Missing product data", http.StatusBadRequest)
		return
	}

	// Формируем тело SOAP-запроса
	payload := fmt.Sprintf("<product><name>%s</name><price>%s</price></product>", name, price)

	// Отправляем SOAP-запрос
	response, err := callSOAPService("add", payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ клиенту
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(response))
}

func getAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Отправляем SOAP-запрос с действием "getAll"
	response, err := callSOAPService("getAll", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ клиенту
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(response))
}

func main() {
	router := mux.NewRouter()

	// REST-эндпоинты
	router.HandleFunc("/api/products", getProductHandler).Methods("GET")
	router.HandleFunc("/api/products", addProductHandler).Methods("POST")
	router.HandleFunc("/api/products/all", getAllProductsHandler).Methods("GET")

	log.Println("REST-to-SOAP Proxy server started on :8083")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
