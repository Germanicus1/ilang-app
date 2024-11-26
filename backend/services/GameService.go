package services

import (
	"backend/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// FetchGames retrieves a list of games from the Supabase database
func FetchGames() ([]Game, error) {
	cfg := config.LoadConfig()

	// Define the Supabase REST API URL for the games table
	url := fmt.Sprintf("%s/rest/v1/games", cfg.SupabaseURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers for Supabase API
	req.Header.Set("apikey", cfg.SupabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.SupabaseKey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch games from Supabase")
	}

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var games []Game
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, err
	}

	return games, nil
}
