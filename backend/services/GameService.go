package services

import (
	"backend/config"
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	req.Header.Set("Prefer", "return=representation") // Ensures Supabase returns the created row

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Game{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return models.Game{}, fmt.Errorf("failed to read Supabase response: %v", err)
	}

	// Log the raw response for debugging
	log.Printf("Response Body: %s\n", string(responseBody))

	// Check for non-201 status codes
	if resp.StatusCode != http.StatusCreated {
		return models.Game{}, fmt.Errorf("failed to create game in Supabase: %s", string(responseBody))
	}

	// Decode the response body into a slice of Game objects
	var createdGames []models.Game
	if len(responseBody) > 0 {
		if err := json.Unmarshal(responseBody, &createdGames); err != nil {
			log.Printf("Error decoding response body: %v\n", err)
			return models.Game{}, fmt.Errorf("failed to decode Supabase response: %v", err)
		}
	}

	// Ensure at least one game is returned
	if len(createdGames) == 0 {
		return models.Game{}, fmt.Errorf("no game created, empty response from Supabase")
	}

	// Return the first game from the array
	return createdGames[0], nil
}
