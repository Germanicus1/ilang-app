package handlers

import (
	"backend/services"
	"encoding/json"
	"log"
	"net/http"
)

type Game struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SubjectID   string `json:"subject_id"`
	Difficulty  int    `json:"difficulty_level"`
	CreatedAt   string `json:"created_at"`
}

// GamesHandler retrieves a list of games from the Supabase database
func GamesHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch games from the Supabase service
	games, err := services.FetchGames()
	if err != nil {
		log.Println("Error fetching games:", err)
		http.Error(w, "Failed to fetch games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}
