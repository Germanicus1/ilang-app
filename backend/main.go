package main

import (
	"log"
	"net/http"

	"backend/config"
	"backend/routes"
	"backend/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// fmt.Println(cfg)
	// Initialize Supabase
	services.InitSupabase(cfg)

	// Set up the HTTP server and routes
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	// Start the server
	log.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
