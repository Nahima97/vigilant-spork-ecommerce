package utils

import (
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
    "net/http"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func HashPassword (password string) (string, error) {
hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
	return "", err
}
return string(hashedPass), nil
}

func ComparePassword (hashedPass, plainPass string) error {
err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
if err != nil {
	return err
}
return nil 
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, status int, msg string) {
    WriteJSON(w, status, map[string]string{"error": msg})
}

func GenerateJWT(secret string, userID int, role string) (string, error) {
    claims := jwt.MapClaims{
        "sub":  userID,
        "role": role,
        "exp":  time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr string, secret string) (int, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil || !token.Valid {
        return 0, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        // sub is float64 in the parsed token
        return int(claims["sub"].(float64)), nil
    }
    return 0, nil
}
