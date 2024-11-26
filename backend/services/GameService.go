package services

import (
	"backend/config"
	"backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// CreateGame creates a new game in the Supabase database
func CreateGame(title, description, subjectID string, difficulty int) (models.Game, error) {
	cfg := config.LoadConfig()

	// Define the Supabase REST API URL for the games table
	url := fmt.Sprintf("%s/rest/v1/games", cfg.SupabaseURL)

	// Prepare the request body
	gameData := map[string]interface{}{
		"title":            title,
		"description":      description,
		"subject_id":       subjectID,
		"difficulty_level": difficulty,
	}
	body, err := json.Marshal(gameData)
	if err != nil {
		return models.Game{}, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return models.Game{}, err
	}

	// Add required headers for Supabase API
	req.Header.Set("apikey", cfg.SupabaseKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.SupabaseKey))
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Game{}, err
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusCreated {
		return models.Game{}, errors.New("failed to create game in Supabase")
	}

	// Parse the response body
	var createdGame models.Game
	if err := json.NewDecoder(resp.Body).Decode(&createdGame); err != nil {
		return models.Game{}, err
	}

	return createdGame, nil
}
