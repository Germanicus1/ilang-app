package handlers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// GamesHandler retrieves a list of games from the Supabase database
func GamesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


	// Fetch games from the Supabase service
	games, err := services.FetchGames(userID)
	if err != nil {
		log.Println("Error fetching games:", err)
		http.Error(w, "Failed to fetch games", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}

// CreateGameHandler creates a new game in the Supabase database
func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to create the game
	game, err := services.CreateGame(req.Title, req.Description, req.SubjectID, req.Difficulty)
	if err != nil {
		log.Println("Error creating game:", err)
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

// GetGameHandler retrieves a single game by its ID from the Supabase database
func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract the game ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Game ID not provided", http.StatusBadRequest)
		return
	}
	gameID := pathParts[2]
	// Call the service to fetch the game
	game, err := services.FetchGameByID(gameID, userID)
	if err != nil {
		log.Printf("Error fetching game: %v\n", err)
		http.Error(w, "Failed to fetch game", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(game)
}

// UpdateGameHandler updates a game by its ID in the Supabase database
func UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the game ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Game ID not provided", http.StatusBadRequest)
		return
	}
	gameID := pathParts[2]

	// Retrieve user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req models.GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to update the game
	updatedGame, err := services.UpdateGameByID(gameID, userID, req)
	if err != nil {
		log.Printf("Error updating game: %v\n", err)
		http.Error(w, "Failed to update game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGame)
}

// DeleteGameHandler deletes a game by its ID from the Supabase database
func DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the game ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Game ID not provided", http.StatusBadRequest)
		return
	}
	gameID := pathParts[2]

	// Retrieve user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the service to delete the game
	err := services.DeleteGameByID(gameID, userID)
	if err != nil {
		log.Printf("Error deleting game: %v\n", err)
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}