package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"vigilant-spork/utils"

	"github.com/golang-jwt/jwt"
    "github.com/gin-gonic/gin"
    "vigilant-spork/repository"

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


func AuthGinMiddleware(userRepo *repository.UserRepo, secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        // Trim the Bearer prefix
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Check if token is blacklisted
        if userRepo.IsTokenBlacklisted(tokenString) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is blacklisted"})
            return
        }

        // Parse and validate the JWT
        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
            }
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        
        c.Set("jwt", token)
        c.Set("jwtTokenString", tokenString) // <-- Add it here

        
        c.Next()
    }
}
