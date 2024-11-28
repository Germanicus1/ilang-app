package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testUrl := "https://selzofrhnkglokruymgh.supabase.co"
	testKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InNlbHpvZnJobmtnbG9rcnV5bWdoIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzIyNzA5ODAsImV4cCI6MjA0Nzg0Njk4MH0.FE0h8ynCGrsgYDbvsjLf3WvuSLH3ukeaCVaMrImn1es"
	testJWT := "EDBEHVaJYYdv4+Zz5mEL7Q8E7UH0Gy+bVNPMyMFWvw7eRu5Ag6IOvYShdX/dKv9VxwXpo2S0cMxe82mnYQVSjA=="
	// os.Setenv("SUPABASE_URL", testUrl)
	// os.Setenv("SUPABASE_KEY", testKey)

	cfg := LoadConfig()

	if cfg.SupabaseURL != testUrl {
		t.Errorf("Expected SupabaseURL to be set, got %v", cfg.SupabaseURL)
	}

	if cfg.SupabaseKey != testKey {
		t.Errorf("Expected SupabaseKey to be set, got %v", cfg.SupabaseKey)
	}

	if cfg.JWTSecret != testJWT {
		t.Errorf("Expected SupabaseKey to be set, got %v", cfg.SupabaseKey)
	}
}
