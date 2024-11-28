package middleware

import (
	"backend/config"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var cfg = config.LoadConfig()
var jwtSecret = cfg.JWTSecret

type UserClaims struct {
	Sub       string   `json:"sub"`       // User ID
	Role      []string `json:"Role"`      // Array of roles
	GivenName string   `json:"GivenName"` // First name
	Surname   string   `json:"Surname"`   // Last name
	Email     string   `json:"Email"`     // User email
	jwt.StandardClaims
}


// AuthMiddleware validates the JWT and injects user claims into the request context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate claims (e.g., expiration)
		if claims.ExpiresAt < time.Now().Unix() {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		// Inject user_id and full claims into the request context
		ctx := context.WithValue(r.Context(), "user_id", claims.Sub) // Only user_id
		ctx = context.WithValue(ctx, "user_claims", claims)         // Full claims
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
