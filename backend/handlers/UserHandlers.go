package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Read and log the raw request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Println("Request Body:", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req CreateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Supabase Auth Signup
	authURL := fmt.Sprintf("%s/auth/v1/signup", os.Getenv("SUPABASE_URL"))
	apiKey := os.Getenv("SUPABASE_KEY")
	payload, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", authURL, bytes.NewBuffer(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "Failed to contact Supabase: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Log the Supabase API response
	bodyBytes, _ = io.ReadAll(resp.Body)
	fmt.Println("Supabase Response:", string(bodyBytes))

	// Parse the Supabase response
	var supabaseResp SupabaseAuthResponse
	err = json.Unmarshal(bodyBytes, &supabaseResp)
	if err != nil {
		http.Error(w, "Failed to decode Supabase response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("Request Method:", r.Method)
	fmt.Println("Request Headers:", r.Header)
	bodyBytes, _ = io.ReadAll(r.Body)
	fmt.Println("Request Body:", string(bodyBytes))

	// Respond with the user details
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(supabaseResp.User)
}
