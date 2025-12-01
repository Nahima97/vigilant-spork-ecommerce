package middleware

import (
	"context"
	"net/http"
	"strings"
	"vigilant-spork/services"
	"vigilant-spork/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"
const JWTTokenKey contextKey = "jwtTokenString"

func AuthMiddleware(secret string, userService *services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorJSON(w, http.StatusUnauthorized, "missing auth header")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			uid, role, err := utils.ValidateJWT(token)
			if err != nil {
				utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
				return
			}

			//storing user and token in context
			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			ctx = context.WithValue(ctx, JWTTokenKey, token)

			isBlacklisted, err := userService.IsTokenBlacklisted(token)
			if err != nil {
				utils.ErrorJSON(w, http.StatusInternalServerError, "error checking token")
				return
			}
			if isBlacklisted {
				utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) int {
	val := ctx.Value(UserIDKey)
	if id, ok := val.(int); ok {
		return id
	}
	return 0
}

func GetUserRole(ctx context.Context) string {
	val := ctx.Value(UserRoleKey)
	if role, ok := val.(string); ok {
		return role
	}
	return ""
}

func GetToken(ctx context.Context) string {
	val := ctx.Value(JWTTokenKey)
	if token, ok := val.(string); ok {
		return token
	}
	return ""
}
