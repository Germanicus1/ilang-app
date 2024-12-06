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

type SupabaseDeleteResponse struct {
	Data  struct {
		User map[string]interface{} `json:"user"`
	} `json:"data"`
	Error *SupabaseError `json:"error"`
}

type SupabaseError struct {
	Message string `json:"message"`
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
    userID := r.PathValue("id")
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    // Query the database for the user
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

type UpdateUserRequest struct {
    Email string `json:"email,omitempty"`
    Role  string `json:"role,omitempty"`
}

func UpdateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.PathValue("id")
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    var updateReq UpdateUserRequest
    err := json.NewDecoder(r.Body).Decode(&updateReq)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if updateReq.Email == "" && updateReq.Role == "" {
        http.Error(w, "No fields to update", http.StatusBadRequest)
        return
    }

    // Prepare the update payload for public.users
    updatePayload := map[string]interface{}{}
    if updateReq.Email != "" {
        updatePayload["email"] = updateReq.Email
    }
    if updateReq.Role != "" {
        updatePayload["role"] = updateReq.Role
    }

    // Update public.users
    var updatedUsers []map[string]interface{}
    err = services.Supabase.DB.From("users").Update(updatePayload).Eq("id", userID).Execute(&updatedUsers)
    if err != nil {
        http.Error(w, "Failed to update user in public.users: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if len(updatedUsers) == 0 {
        http.Error(w, "User not found in public.users", http.StatusNotFound)
        return
    }

    // Update auth.users if email has changed
    if updateReq.Email != "" {
        authPayload := map[string]interface{}{
            "email": updateReq.Email,
        }

        authURL := fmt.Sprintf("%s/auth/v1/admin/users/%s", os.Getenv("SUPABASE_URL"), userID)
        reqBody, _ := json.Marshal(authPayload)

        authReq, _ := http.NewRequest("PUT", authURL, bytes.NewBuffer(reqBody))
        authReq.Header.Set("Content-Type", "application/json")
        authReq.Header.Set("apikey", os.Getenv("SUPABASE_KEY"))
        authReq.Header.Set("Authorization", "Bearer "+os.Getenv("SERVICE_ROLE_KEY"))

        client := &http.Client{}
        resp, err := client.Do(authReq)
        if err != nil || resp.StatusCode != http.StatusOK {
            http.Error(w, "Failed to update user in auth.users: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Respond with the updated user
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedUsers[0])
}

func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Step 1: Delete from public.users
	var deletedUsers []map[string]interface{}
	err := services.Supabase.DB.From("users").Delete().Eq("id", userID).Execute(&deletedUsers)
	if err != nil {
		http.Error(w, "Failed to delete user from public.users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// if len(deletedUsers) == 0 {
	// 	http.Error(w, "User not found in public.users", http.StatusNotFound)
	// 	return
	// }

	// Step 2: Delete from auth.users using Supabase Admin API
	adminDeleteUser(w, userID)
}

func adminDeleteUser(w http.ResponseWriter, userID string) {
	authURL := fmt.Sprintf("%s/auth/v1/admin/users/%s", os.Getenv("SUPABASE_URL"), userID)

	authReq, _ := http.NewRequest("DELETE", authURL, nil)
	authReq.Header.Set("Authorization", "Bearer "+os.Getenv("SERVICE_ROLE_KEY"))
	authReq.Header.Set("apikey", os.Getenv("SERVICE_ROLE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(authReq)
	if err != nil {
		http.Error(w, "Failed to delete user from auth.users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Log for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Auth API Response Body: %s\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		http.Error(w, fmt.Sprintf("Failed to delete user from auth.users: %s", resp.Status), resp.StatusCode)
		return
	}

		// Parse the response body
	var supabaseResp SupabaseDeleteResponse
	err = json.Unmarshal(bodyBytes, &supabaseResp)
	if err != nil {
		http.Error(w, "Failed to parse response from auth.users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for errors in the response
	if supabaseResp.Error != nil {
		http.Error(w, "Supabase error: "+supabaseResp.Error.Message, http.StatusInternalServerError)
		return
	}

	// Log the deleted user's information
	fmt.Printf("Deleted User: %+v\n", supabaseResp.Data.User)
}
