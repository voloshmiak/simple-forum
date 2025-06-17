package auth

import (
	"forum-project/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JwtAuthenticator struct {
	secret      string
	expiryHours int
}

func NewJwtAuthenticator(secret string, expiryHours int) *JwtAuthenticator {
	return &JwtAuthenticator{
		secret:      secret,
		expiryHours: expiryHours,
	}
}

func (a *JwtAuthenticator) GenerateToken(user *model.User) (string, error) {
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

func (a *JwtAuthenticator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func (a *JwtAuthenticator) GetClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
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
