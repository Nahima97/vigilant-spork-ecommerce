package middleware

import (
	"context"
	"net/http"
	"strings"
	"vigilant-spork/utils"
)

type contextKey string

var userIDKey contextKey = "userID"
var userRoleKey contextKey = "userRole"

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

			ctx := context.WithValue(r.Context(), userIDKey, uid)
			ctx = context.WithValue(ctx, userRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) int {
	if id, ok := ctx.Value(userIDKey).(int); ok {
		return id
	}
	return 0
}

func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(userRoleKey).(string); ok {
		return role
	}
	return ""
}
