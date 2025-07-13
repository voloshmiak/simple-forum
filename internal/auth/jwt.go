package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var (
	ErrZeroID     = errors.New("id cannot be 0")
	ErrEmptyName  = errors.New("username cannot be empty")
	ErrEmptyRole  = errors.New("role cannot be empty")
	ErrNilRequest = errors.New("request cannot be nil")
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

func (a *JWTAuthenticator) GenerateToken(userID int, userName, userRole string) (string, error) {
	if userID == 0 {
		return "", ErrZeroID
	}

	if userName == "" {
		return "", ErrEmptyName
	}

	if userRole == "" {
		return "", ErrEmptyRole
	}

	claims := jwt.MapClaims{
		"user": map[string]interface{}{
			"id":   userID,
			"name": userName,
			"role": userRole,
		},
		"exp": time.Now().Add(time.Hour * time.Duration(a.expiryHours)).Unix(),
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
	if r == nil {
		return nil, ErrNilRequest
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	return a.ValidateToken(cookie.Value)
}
