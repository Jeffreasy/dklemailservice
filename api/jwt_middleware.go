package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// HandleJWTMiddleware is middleware voor JWT authenticatie
func HandleJWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Haal JWT uit Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header ontbreekt", http.StatusUnauthorized)
			return
		}

		// Check of Authorization header begint met "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Ongeldige Authorization header", http.StatusUnauthorized)
			return
		}

		// Haal token uit header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse en valideer token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Ongeldige token", http.StatusUnauthorized)
			return
		}

		// Ga verder met de request
		next(w, r)
	}
}
