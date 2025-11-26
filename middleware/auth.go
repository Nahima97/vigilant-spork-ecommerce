package middleware

import (
	"context"
	"net/http"
	"strings"
	"vigilant-spork/utils"

	"github.com/gofrs/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

// AuthMiddleware validates JWT and stores user ID (as uuid.UUID) and role in context
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorJSON(w, http.StatusUnauthorized, "missing auth header")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userIDStr, role, err := utils.ValidateJWT(token)
			if err != nil {
				utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
				return
			}

			// Convert userID string to uuid.UUID
			userID, err := uuid.FromString(userIDStr)
			if err != nil {
				utils.ErrorJSON(w, http.StatusUnauthorized, "invalid user ID in token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves uuid.UUID from context
func GetUserID(ctx context.Context) uuid.UUID {
	val := ctx.Value(UserIDKey)
	if id, ok := val.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

// GetUserRole retrieves role string from context
func GetUserRole(ctx context.Context) string {
	val := ctx.Value(UserRoleKey)
	if role, ok := val.(string); ok {
		return role
	}
	return ""
}
