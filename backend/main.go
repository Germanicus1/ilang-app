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

	log.Printf("Supabase URL: %s\n", cfg.SupabaseURL)
	log.Printf("Supabase Key: %s\n", cfg.SupabaseKey)

	// Register routes
	// http.HandleFunc("/health", handlers.HealthHandler)
	// http.HandleFunc("/games", handlers.GamesHandler)         // GET all games
	http.HandleFunc("/games", handlers.CreateGameHandler)    // POST a game
	http.HandleFunc("/games/", handlers.GetGameHandler)      // GET /games/{id}

	// Start the server
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}