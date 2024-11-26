package middleware

import (
	"log"
	"net/http"
	"os"
)

var errorLogFile *os.File

func InitErrorLogger() {
	var err error
	errorLogFile, err = os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to initialize error logger: %v", err)
	}
	log.SetOutput(errorLogFile)
}

func ErrorLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Error: %v\n", err)
				http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}