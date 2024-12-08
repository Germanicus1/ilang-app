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
	"time"
)

// Common structs for requests and responses
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}

// Utility functions for CORS and request parsing
func handlePreflight(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func parseRequestBody[T any](r *http.Request) (T, error) {
	var req T
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return req, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return req, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return req, nil
}

// Supabase API helpers
func callSupabaseAPI(method, url string, payload interface{}, headers map[string]string) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// CreateUserHandler
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Handle preflight (CORS)
    if r.Method == http.MethodOptions {
        handlePreflight(w)
        return
    }

    // Parse request payload
    req, err := parseRequestBody[CreateUserRequest](r)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Prepare Supabase API call
    authURL := fmt.Sprintf("%s/auth/v1/signup", os.Getenv("SUPABASE_URL"))
    headers := map[string]string{
        "Content-Type": "application/json",
        "apikey":       os.Getenv("SUPABASE_KEY"),
    }

    // Call Supabase API
    resp, err := callSupabaseAPI(http.MethodPost, authURL, req, headers)
    if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated) {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read the entire response body once
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read response body", http.StatusInternalServerError)
        return
    }

    // Parse the response body
    var supabaseResp struct {
        AccessToken  string       `json:"access_token"`
        TokenType    string       `json:"token_type"`
        ExpiresIn    int          `json:"expires_in"`
        ExpiresAt    int64        `json:"expires_at"`
        RefreshToken string       `json:"refresh_token"`
        User         SupabaseUser `json:"user"`
    }

    err = json.Unmarshal(bodyBytes, &supabaseResp)
    if err != nil {
        http.Error(w, "Failed to parse response", http.StatusInternalServerError)
        return
    }

    // Respond with user details
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(supabaseResp.User)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var users []models.User
	err := services.Supabase.DB.From("users").Select("*").Eq("id", userID).Execute(&users)
	if err != nil || len(users) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users[0])
}

func UpdateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	updateReq, err := parseRequestBody[UpdateUserRequest](r)
	if err != nil || (updateReq.Email == "" && updateReq.Role == "") {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
		return
	}

	updatePayload := map[string]interface{}{}
	if updateReq.Email != "" {
		updatePayload["email"] = updateReq.Email
	}
	if updateReq.Role != "" {
		updatePayload["role"] = updateReq.Role
	}

	var updatedUsers []map[string]interface{}
	err = services.Supabase.DB.From("users").Update(updatePayload).Eq("id", userID).Execute(&updatedUsers)
	if err != nil || len(updatedUsers) == 0 {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	if updateReq.Email != "" {
		authURL := fmt.Sprintf("%s/auth/v1/admin/users/%s", os.Getenv("SUPABASE_URL"), userID)
		headers := map[string]string{
			"Content-Type":  "application/json",
			"apikey":        os.Getenv("SUPABASE_KEY"),
			"Authorization": "Bearer " + os.Getenv("SERVICE_ROLE_KEY"),
		}
		_, err := callSupabaseAPI(http.MethodPut, authURL, map[string]string{"email": updateReq.Email}, headers)
		if err != nil {
			http.Error(w, "Failed to update auth.users", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUsers[0])
}

func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the {id} from the URL path
	userID := r.URL.Path[len("/users/"):] // Extract the part of the path after "/users/"
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Step 1: Delete from public.users
	err := services.Supabase.DB.From("users").Delete().Eq("id", userID).Execute(nil)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from auth.users using Supabase Admin API
	authURL := fmt.Sprintf("%s/auth/v1/admin/users/%s", os.Getenv("SUPABASE_URL"), userID)
	headers := map[string]string{
		"apikey":        os.Getenv("SERVICE_ROLE_KEY"),
		"Authorization": "Bearer " + os.Getenv("SERVICE_ROLE_KEY"),
	}
	_, err = callSupabaseAPI(http.MethodDelete, authURL, nil, headers)
	if err != nil {
		http.Error(w, "Failed to delete auth user", http.StatusInternalServerError)
		return
	}

	// Respond with No Content status
	w.WriteHeader(http.StatusNoContent)
}
