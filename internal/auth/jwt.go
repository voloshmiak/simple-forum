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

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *JWTAuthenticator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func (a *JWTAuthenticator) GetClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	claims, err := a.ValidateToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
