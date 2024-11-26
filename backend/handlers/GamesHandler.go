package handlers

import (
	"backend/services"
	"encoding/json"
	"log"
	"net/http"
)

type CreateGameRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SubjectID   string `json:"subject_id"`
	Difficulty  int    `json:"difficulty_level"`
}

// CreateGameHandler creates a new game in the Supabase database
func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req CreateGameRequest
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
