package services

import (
	"backend/config"
	"log"

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
