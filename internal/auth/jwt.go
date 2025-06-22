package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"simple-forum/internal/model"
	"time"
)

type JWTAuthenticator struct {
	secret      string
	expiryHours int
}

func NewJWTAuthenticator(secret string, expiryHours int) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:      secret,
		expiryHours: expiryHours,
	}
}

func (a *JWTAuthenticator) GenerateToken(user *model.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	authorizedUser := model.AuthorizedUser{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	claims := jwt.MapClaims{
		"user": authorizedUser,
		"exp":  time.Now().Add(time.Hour * time.Duration(a.expiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.secret))

	return signedToken, err
}

func (a *JWTAuthenticator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func (a *JWTAuthenticator) GetClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	return a.ValidateToken(cookie.Value)
}
