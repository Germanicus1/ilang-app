package main

import (
	"log"
	"net/http"

	"backend/config"
	"backend/middleware"
	"backend/routes"
	"backend/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	services.InitSupabase(cfg)

	// Create a base multiplexer
	mux := http.NewServeMux()
	// Register public routes (no middleware)
	routes.RegisterPublicRoutes(mux)

	// Create a sub-mux for secured routes
	securedMux := http.NewServeMux()
	routes.RegisterSecuredRoutes(securedMux)

	// Wrap the secured mux in middleware
	mux.Handle("/users/", middleware.ValidateJWT(securedMux))

	// Start the server
	log.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
