package services

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) CreateUser(email, password string) (models.SupabaseUser, error) {
    authURL := fmt.Sprintf("%s/auth/v1/signup", os.Getenv("SUPABASE_URL"))
    headers := map[string]string{
        "Content-Type": "application/json",
        "apikey":       os.Getenv("SUPABASE_KEY"),
    }

    payload := map[string]string{"email": email, "password": password}
    resp, err := callSupabaseAPI(http.MethodPost, authURL, payload, headers)
    if err != nil {
        return models.SupabaseUser{}, err
    }
    defer resp.Body.Close()

    var supabaseResp struct {
        User models.SupabaseUser `json:"user"`
    }
    json.NewDecoder(resp.Body).Decode(&supabaseResp)
    return supabaseResp.User, nil
}
