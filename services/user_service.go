package services

import (
	"errors"
	"os"
	"regexp"
	"vigilant-spork/middleware"
	"vigilant-spork/models"
	"vigilant-spork/repository"
	"vigilant-spork/utils"
)

type UserService struct {
	UserRepo repository.UserRepository
}

var ErrEmailExists = errors.New("email already registered")

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func (s *UserService) RegisterUser(user *models.User) error {
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	existing, err := s.UserRepo.GetUserByEmail(user.Email)
	if err == nil && existing != nil {
		return ErrEmailExists
	}

	hashedPass, err := utils.HashPassword(user.Password) 
	if err != nil {
		return err
	}
	user.Password = hashedPass

	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
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

	token, err := middleware.GenerateJWT(secret, user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) AddTokenToBlacklist(token string) error {
	err := s.UserRepo.AddTokenToBlacklist(token)
	if err != nil {
		return err
	}
	return nil
}
