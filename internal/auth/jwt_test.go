package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuthenticator_GenerateToken(t *testing.T) {
	t.Parallel()
	authenticator := NewJWTAuthenticator("mysecretkey", 24)

	tests := []struct {
		name     string
		userID   int
		userName string
		userRole string
		valid    bool
		err      error
	}{
		{
			name:     "Valid User",
			userID:   1,
			userName: "testuser",
			userRole: "user",
			valid:    true,
		},
		{
			name:     "Invalid User ID",
			userID:   0,
			userName: "testuser",
			userRole: "user",
			valid:    false,
			err:      ErrZeroID,
		},
		{
			name:     "Invalid User name",
			userID:   1,
			userName: "",
			userRole: "user",
			valid:    false,
			err:      ErrEmptyName,
		},
		{
			name:     "Invalid User role",
			userID:   1,
			userName: "testuser",
			userRole: "",
			valid:    false,
			err:      ErrEmptyRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authenticator.GenerateToken(tt.userID, tt.userName, tt.userRole)

			if tt.valid {
				if err != nil {
					t.Errorf("%s: no error expected, but got %s", tt.name, err.Error())
				}
				if token == "" {
					t.Errorf("%s: expected a token, got an empty string", tt.name)
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected %s, got nil", tt.name, tt.err)
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err.Error(), err.Error())
				}
				if token != "" {
					t.Errorf("%s: expected no token, but got %s", tt.name, token)
				}
			}
		})
	}
}

func generateTestToken(secret string, expiryHours int, userID int, userName, userRole string) string {
	claims := jwt.MapClaims{
		"user": map[string]interface{}{
			"id":   userID,
			"name": userName,
			"role": userRole,
		},
		"exp": time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func TestJWTAuthenticator_ValidateToken(t *testing.T) {
	t.Parallel()
	authenticator := NewJWTAuthenticator("mysecretkey", 24)
	wrongSecretAuthenticator := NewJWTAuthenticator("wrongsecret", 24)

	tests := []struct {
		name           string
		token          string
		expectedClaims jwt.MapClaims
		valid          bool
		err            error
	}{
		{
			name:  "Valid Token",
			token: generateTestToken(authenticator.secret, 1, 1, "testuser", "admin"),
			expectedClaims: jwt.MapClaims{
				"user": map[string]interface{}{
					"id":   float64(1),
					"name": "testuser",
					"role": "admin",
				},
			},
			valid: true,
		},
		{
			name:  "Malformed Token",
			token: "token",
			valid: false,
			err:   jwt.ErrTokenMalformed,
		},
		{
			name:  "Expired Token",
			token: generateTestToken(authenticator.secret, -1, 1, "testuser", "admin"),
			valid: false,
			err:   jwt.ErrTokenExpired,
		},
		{
			name:  "Wrong Signature Token",
			token: generateTestToken(wrongSecretAuthenticator.secret, 1, 1, "testuser", "admin"),
			valid: false,
			err:   jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authenticator.ValidateToken(tt.token)

			if tt.valid {
				if err != nil {
					t.Errorf("%s: no error expected, but got %s", tt.name, err.Error())
				}
				if claims == nil {
					t.Errorf("%s: expected claims, got nil", tt.name)
				}
				if !reflect.DeepEqual(claims["user"], tt.expectedClaims["user"]) {
					t.Errorf("%s: expected claims %s, got %s", tt.name, tt.expectedClaims["user"], claims["user"])
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected %s, got nil", tt.name, tt.err)
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err, err.Error())
				}
				if claims != nil {
					t.Errorf("%s: expected no claims, but got %s", tt.name, claims)
				}
			}
		})
	}
}

func TestJWTAuthenticator_GetClaimsFromRequest(t *testing.T) {
	t.Parallel()
	authenticator := NewJWTAuthenticator("mysecretkey", 24)

	tests := []struct {
		name           string
		request        *http.Request
		expectedClaims jwt.MapClaims
		valid          bool
		err            error
	}{
		{
			name: "Valid Request with Cookie",
			request: func() *http.Request {
				token := generateTestToken(authenticator.secret, 1, 1, "testuser", "admin")
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{Name: "token", Value: token})
				return r
			}(),
			expectedClaims: jwt.MapClaims{
				"user": map[string]interface{}{
					"id":   float64(1),
					"name": "testuser",
					"role": "admin",
				},
			},
			valid: true,
		},
		{
			name: "Request with Invalid Token",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{Name: "token", Value: "invalidtoken"})
				return r
			}(),
			valid: false,
			err:   jwt.ErrTokenMalformed,
		},
		{
			name: "Request without Token",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				return r
			}(),
			valid: false,
			err:   http.ErrNoCookie,
		},
		{
			name:    "Nil Request",
			request: nil,
			valid:   false,
			err:     ErrNilRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authenticator.GetClaimsFromRequest(tt.request)

			if tt.valid {
				if err != nil {
					t.Errorf("%s: no error expected, but got %s", tt.name, err.Error())
				}
				if claims == nil {
					t.Errorf("%s: expected claims, got nil", tt.name)
				}
				if !reflect.DeepEqual(claims["user"], tt.expectedClaims["user"]) {
					t.Errorf("%s: expected claims %s, got %s", tt.name, tt.expectedClaims["user"], claims["user"])
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected an error, got nil", tt.name)
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err.Error(), err.Error())
				}
				if claims != nil {
					t.Errorf("%s: expected no claims, but got %s", tt.name, claims)
				}
			}
		})
	}
}
