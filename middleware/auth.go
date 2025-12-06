package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"vigilant-spork/db"
	"vigilant-spork/models"
	"vigilant-spork/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"
const JWTTokenKey contextKey = "jwtTokenString"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorJSON(w, http.StatusUnauthorized, "missing auth header")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			uid, role, err := ValidateJWT(token, secret)
			if err != nil {
				utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
				return
			}

			userUUID, err := uuid.FromString(uid)
			if err != nil {
				utils.ErrorJSON(w, http.StatusUnauthorized, "invalid user ID in token")
				return
			}

			isBlacklisted, err := IsTokenBlacklisted(token)
			if err != nil {
				utils.ErrorJSON(w, http.StatusInternalServerError, "error checking token")
				return
			}

			if isBlacklisted {
				utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userUUID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			ctx = context.WithValue(ctx, JWTTokenKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GenerateJWT(secret string, userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr string, secret string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		uid, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)

		return uid, role, nil
	}

	return "", "", fmt.Errorf("invalid sub claim")
}

func IsTokenBlacklisted(token string) (bool, error) {
	var entry models.BlacklistedToken
	err := db.Db.Where("token = ?", token).First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetUserID(ctx context.Context) uuid.UUID {
	val := ctx.Value(UserIDKey)
	if id, ok := val.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
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
