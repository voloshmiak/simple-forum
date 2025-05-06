package auth

import (
	"errors"
	"forum-project/internal/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret-key"))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return []byte("secret-key"), nil })
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func ValidateTokenFromRequest(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	token, err := ValidateToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return token, nil
}
