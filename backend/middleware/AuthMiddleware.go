package middleware

import (
	"backend/config"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	UserIDContextKey contextKey = "userID"
	RoleContextKey   contextKey = "role"
)

// Middleware to validate JWT
func ValidateJWT(next http.Handler) http.Handler {
	cfg := config.LoadConfig() // loading .env

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := cfg.JWTSecret
			if secret == "" {
				return nil, fmt.Errorf("JWT secret is missing")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			http.Error(w, "Invalid token: 'sub' claim is missing", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			http.Error(w, "Invalid token: 'role' claim is missing", http.StatusUnauthorized)
			return
		}

		// Add userID and role to request context
		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		ctx = context.WithValue(ctx, RoleContextKey, role)

		// Proceed to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
