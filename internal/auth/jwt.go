package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var (
	ZeroIDErr     = errors.New("id cannot be 0")
	EmptyNameErr  = errors.New("username cannot be empty")
	EmptyRoleErr  = errors.New("role cannot be empty")
	NilRequestErr = errors.New("request cannot be nil")
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

func (a *JWTAuthenticator) GenerateToken(id int, name, role string) (string, error) {
	if id == 0 {
		return "", ZeroIDErr
	}

	if name == "" {
		return "", EmptyNameErr
	}

	if role == "" {
		return "", EmptyRoleErr
	}

	user := map[string]interface{}{
		"id":   id,
		"name": name,
		"role": role,
	}

	claims := jwt.MapClaims{
		"user": user,
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
	if r == nil {
		return nil, NilRequestErr
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	return a.ValidateToken(cookie.Value)
}
