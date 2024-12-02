package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting catalog service...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
		w.Write([]byte("Hello, World!"))
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}

	log.Println("Server stopped")
}
