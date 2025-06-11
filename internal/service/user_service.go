package service

import (
	"errors"
	"forum-project/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrMismatchPassword       = errors.New("passwords do not match")
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
	ErrUserNameAlreadyExists  = errors.New("username already exists")
	ErrInvalidToken           = errors.New("invalid token")
)

type UserStorage interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	InsertUser(user *model.User) (int, error)
}

type UserService struct {
	repository  UserStorage
	jwtSecret   string
	expiryHours int
}

func NewUserService(repository UserStorage, jwtSecret string, expiryHours int) *UserService {
	return &UserService{
		repository:  repository,
		jwtSecret:   jwtSecret,
		expiryHours: expiryHours,
	}
}

func (u *UserService) Login(email, password string) (string, error) {
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

	authorizedUser := model.AuthorizedUser{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	claims := jwt.MapClaims{
		"user": authorizedUser,
		"exp":  time.Now().Add(time.Hour * time.Duration(u.expiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(u.jwtSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
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
		Username: username,
		Email:    email, Role: "user",
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

func ValidateToken(tokenString, jwtSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
