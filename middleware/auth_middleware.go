package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(jwtsecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, `{"error": "authorization header required"}`, http.StatusUnauthorized)

				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeJSONError(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)

				return
			}

			tokenString := parts[1]
			token, tErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtsecret), nil
			})

			if tErr != nil || !token.Valid {
				writeJSONError(w, `{"error": "invalid token"}`, http.StatusUnauthorized)

				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeJSONError(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)

				return
			}
			userID := uint(claims["user_id"].(float64))

			// Set user ID in context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// Continue to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": ` + message + `}`))
}
func GetUserIDFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value(UserIDKey).(uint)
	if !ok {
		return 0
	}
	return userID
}
