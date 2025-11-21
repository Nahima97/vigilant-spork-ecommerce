package middleware

import (
    "context"
    "net/http"
    "strings"
	"vigilant-spork/utils"
)

type contextKey string
const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "UserRole"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                utils.ErrorJSON(w, http.StatusUnauthorized, "missing auth header")
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            uid, role, err := utils.ValidateJWT(token, secret)
            if err != nil {
                utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
                return
            }

            ctx := context.WithValue(r.Context(), UserIDKey, uid)
			ctx = context.WithValue(ctx, UserRoleKey, role)
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
