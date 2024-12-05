package handlers

import (
	"backend/models"
	"backend/services"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Structs for request and response
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type SupabaseAuthResponse struct {
	User SupabaseUser `json:"user"`
}

// CreateUserHandler handles user creation
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight requests (CORS)
	if r.Method == http.MethodOptions {
		handlePreflight(w)
		return
	}

	// Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	req, err := parseRequestBody(r)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call Supabase to create the user
	user, err := createSupabaseUser(req)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with user details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// handlePreflight handles CORS preflight requests
func handlePreflight(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

// parseRequestBody parses the request body into CreateUserRequest
func parseRequestBody(r *http.Request) (CreateUserRequest, error) {
	var req CreateUserRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return req, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	fmt.Println("Request Body:", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		return req, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return req, nil
}

// createSupabaseUser creates a user in Supabase Auth
func createSupabaseUser(req CreateUserRequest) (SupabaseUser, error) {
	authURL := fmt.Sprintf("%s/auth/v1/signup", os.Getenv("SUPABASE_URL"))
	apiKey := os.Getenv("SUPABASE_KEY")

	payload, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", authURL, bytes.NewBuffer(payload))
	if err != nil {
		return SupabaseUser{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return SupabaseUser{}, fmt.Errorf("failed to contact Supabase: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Supabase Response:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return SupabaseUser{}, fmt.Errorf("supabase error: %s", resp.Status)
	}

	var supabaseResp SupabaseAuthResponse
	err = json.Unmarshal(bodyBytes, &supabaseResp)
	if err != nil {
		return SupabaseUser{}, fmt.Errorf("failed to parse Supabase response: %w", err)
	}

	return supabaseResp.User, nil
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Extract userID from the URL path
    userID := r.URL.Query().Get("id")
    if userID == "" {
        http.Error(w, "Missing user ID", http.StatusBadRequest)
        return
    }

    var users []models.User
    err := services.Supabase.DB.From("users").Select("*").Eq("id", userID).Execute(&users)
    if err != nil {
        http.Error(w, "Failed to retrieve user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if len(users) == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Return the user as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users[0])
}