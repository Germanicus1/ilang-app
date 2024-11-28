package main

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"backend/services"
	"log"
	"net/http"
)

func main() {
	// Load configuration and iniot Supabase client
	cfg := config.LoadConfig()
	services.InitSupabase(cfg)

	// Log Supabase configuration for debugging
	log.Printf("Supabase URL: %s\n", cfg.SupabaseURL)
	log.Printf("Supabase client initializec")
	log.Printf("JWT_SECRET: %s\n", cfg.JWTSecret)

	// Register routes
	http.HandleFunc("/health", handlers.HealthHandler) // Health check

	// Handle /games for both GET and POST methods
	http.HandleFunc("/games", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GamesHandler(w, r) // Handle GET /games
		case http.MethodPost:
			handlers.CreateGameHandler(w, r) // Handle POST /games
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Handle /games/{id}
	http.HandleFunc("/games/", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGameHandler(w, r) // GET /games/{id}
		case http.MethodPut:
			handlers.UpdateGameHandler(w, r) // PUT /games/{id}
		case http.MethodDelete:
			handlers.DeleteGameHandler(w, r) // DELETE /games/{id}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Start the server
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
