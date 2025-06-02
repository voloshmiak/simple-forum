package service

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")
var ErrMissmatchPassword = errors.New("passwords do not match")
var ErrUserEmailAlreadyExists = errors.New("user email already exists")
var ErrUserNameAlreadyExists = errors.New("username already exists")

type UserStorage interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	InsertUser(user *model.User) (int, error)
}

type UserService struct {
	repository UserStorage
}

func NewUserService(repository UserStorage) *UserService {
	return &UserService{repository: repository}
}

func (u *UserService) Authenticate(email, password, jwtSecret string, expiryHours int) (string, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrWrongPassword
	}

	token, err := auth.GenerateToken(user, jwtSecret, expiryHours)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) Register(username, email, password1, password2 string) error {
	if password1 != password2 {
		return ErrMissmatchPassword
	}

	existingUserByUsername, err := u.repository.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existingUserByUsername != nil {
		return ErrUserNameAlreadyExists
	}

	existingUserByEmail, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if existingUserByEmail != nil {
		return ErrUserEmailAlreadyExists
	}

	user := &model.User{Username: username, Email: email, Role: "user", CreatedAt: time.Now()}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	_, err = u.repository.InsertUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserByID(id int) (*model.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
