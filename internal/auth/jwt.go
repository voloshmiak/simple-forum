package auth

import (
	"errors"
	"forum-project/internal/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *model.User, jwtSecret string, expiryHours int) (string, error) {
	authorizedUser := model.AuthorizedUser{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	claims := jwt.MapClaims{
		"user": authorizedUser,
		"exp":  time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func GetClaimsFromRequest(r *http.Request, jwtSecret string) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	token, err := ValidateToken(cookie.Value, jwtSecret)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
