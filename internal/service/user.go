package service

import (
	"errors"
	"forum-project/internal/models"
	"forum-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
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

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (u *UserService) Register(username, email, password1, password2 string) (*models.User, error) {
	if password1 != password2 {
		return nil, errors.New("passwords do not match")
	}

	user := models.NewUser()
	user.Username = username
	user.Email = email

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)

	user.PasswordHash = string(hashedPassword)
	user.CreatedAt = "Now"
	user.Role = "user"

	userid, err := u.repository.InsertUser(user)
	if err != nil {
		return nil, err
	}

	user.ID = userid
	return user, nil
}

func (u *UserService) GetUserByMail(email string) (*models.User, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
