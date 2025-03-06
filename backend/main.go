package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables
	// Create a new router
	router := mux.NewRouter()

	// Define a test endpoint
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong! Server is running 🚀"))
	}).Methods("GET")

	// Start the server
	port := ":8080"
	fmt.Println("🚀 Backend server running on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
