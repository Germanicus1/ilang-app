package services

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nedpals/supabase-go"
)

var Supabase *supabase.Client

// InitSupabase initializes the Supabase client
func InitSupabase(cfg config.Config) {
	supabaseURL := cfg.SupabaseURL
	supabaseKey := cfg.SupabaseKey

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatalf("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	client := supabase.CreateClient(supabaseURL, supabaseKey)
	Supabase = client
}


// callSupabaseAPI sends an HTTP request to the Supabase API and returns the response
func callSupabaseAPI(method, url string, payload interface{}, headers map[string]string) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		// Marshal the payload to JSON
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(data)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}
