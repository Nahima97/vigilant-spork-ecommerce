package services

import (
	"vigilant-spork/models"
	"vigilant-spork/repository"
	"vigilant-spork/utils"
)

type UserService struct {
	UserRepo repository.UserRepository
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

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}


