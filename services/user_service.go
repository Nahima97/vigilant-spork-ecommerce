package services

import (
	"errors"
	"fmt"
	"regexp"
	"vigilant-spork/models"
	"vigilant-spork/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	
)

type UserService struct {
	UserRepo repository.UserRepository
	secret string
}



// error for duplicate email
var ErrEmailExists = errors.New("email already registered")

func (s *UserService) RegisterUser(user *models.User) error {
	// Validate email
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Check if email exists
	existing, err := s.UserRepo.GetUserByEmail(user.Email)
	if err == nil && existing != nil {
		return ErrEmailExists
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	// Default role
	if user.Role == "" {
		user.Role = "customer"
	}

	// Save user
	return s.UserRepo.CreateUser(user)
}

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func (s *UserService) RevokeSession(tokenString string) error {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(s.secret), nil
    })

    if err != nil || !token.Valid {
        return errors.New("invalid or expired token")
    }

    return s.UserRepo.AddTokenToBlacklist(tokenString)
}