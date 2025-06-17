package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"simple-forum/internal/model"
	"time"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrMismatchPassword       = errors.New("passwords do not match")
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
	ErrUserNameAlreadyExists  = errors.New("username already exists")
)

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
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) Login(email, password string) (*model.User, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrWrongPassword
	}

	return user, nil
}

func (u *UserService) Register(username, email, password1, password2 string) error {
	if password1 != password2 {
		return ErrMismatchPassword
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

	user := &model.User{
		Username:  username,
		Email:     email,
		Role:      "user",
		CreatedAt: time.Now(),
	}

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
