package services

import (
	"errors"
	"os"
	"regexp"
	"vigilant-spork/models"
	"vigilant-spork/repository"
	"vigilant-spork/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repository.UserRepository
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

func (s *UserService) Login(login *models.User) (string, error) {
	
	user, err := s.UserRepo.GetUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	err = utils.ComparePassword(user.Password, login.Password)
	if err != nil {
		return "", err
	}

	var secret = os.Getenv("JWT_SECRET")

	token, err := utils.GenerateJWT(secret, user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
