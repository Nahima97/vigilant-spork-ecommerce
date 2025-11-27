package middleware

import (
	"context"
	"net/http"
	"strings"
	"vigilant-spork/utils"
)

type contextKey string

var UserIDKey contextKey = "userID"
var UserRoleKey contextKey = "userRole"

// AuthMiddleware validates JWT and stores user ID (as uuid.UUID) and role in context
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorJSON(w, http.StatusUnauthorized, "missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.ErrorJSON(w, http.StatusUnauthorized, "invalid Authorization format")
				return
			}

			token := parts[1]
			uid, role, err := utils.ValidateJWT(token, secret)
			if err != nil {
				utils.ErrorJSON(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) int {
	if id, ok := ctx.Value(UserIDKey).(int); ok {
		return id
	}
	return 0
}

// GetUserRole retrieves role string from context
func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}
