package services

import (
	"backend/config"

	"github.com/nedpals/supabase-go"
)

var Supabase *supabase.Client

func InitSupabase(cfg config.Config) {
	Supabase = supabase.CreateClient(cfg.SupabaseURL, cfg.SupabaseKey)
}