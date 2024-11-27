package main

import (
	"backend/config"
	"backend/handlers"
	"log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Log Supabase configuration for debugging
	log.Printf("Supabase URL: %s\n", cfg.SupabaseURL)
	log.Printf("Supabase Key: %s\n", cfg.SupabaseKey)

	// Register routes
	http.HandleFunc("/health", handlers.HealthHandler) // Health check

	// Handle /games for both GET and POST methods
	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GamesHandler(w, r) // Handle GET /games
		case http.MethodPost:
			handlers.CreateGameHandler(w, r) // Handle POST /games
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle /games/{id}
	http.HandleFunc("/games/", handlers.GetGameHandler)

	// Start the server
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
