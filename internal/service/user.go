package service

import (
	"errors"
	"forum-project/internal/models"
	"forum-project/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (u *UserService) Authenticate(email, password string) (*models.User, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.PasswordHash != password {
		return nil, errors.New("wrong password")
	}

	return user, nil
}
