package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    SupabaseURL string
    SupabaseKey string
    JWTSecret string
}

func LoadConfig() Config {
    err := godotenv.Load("../../backend/backend.env")
    if err != nil {
        log.Fatal("Error loading .env file", err)
    }

    return Config{
        SupabaseURL: os.Getenv("SUPABASE_URL"),
        SupabaseKey: os.Getenv("SUPABASE_KEY"),
        JWTSecret: os.Getenv("JWT_SECRET"),
    }
}
