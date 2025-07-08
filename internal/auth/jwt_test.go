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

func TestTokenPayload_Valid(t *testing.T) {
	tests := []struct {
		name    string
		payload *TokenPayload
		valid   bool
		err     error
	}{
		{
			name: "Valid User",
			payload: &TokenPayload{
				ID:   1,
				Name: "testuser",
				Role: "user",
			},
			valid: true,
		},
		{
			name: "Zero User ID",
			payload: &TokenPayload{
				ID:   0,
				Name: "testuser",
				Role: "user",
			},
			valid: false,
			err:   ErrZeroID,
		},
		{
			name: "Empty User name",
			payload: &TokenPayload{
				ID:   1,
				Name: "",
				Role: "user",
			},
			valid: false,
			err:   ErrEmptyName,
		},
		{
			name: "Empty User role",
			payload: &TokenPayload{
				ID:   1,
				Name: "testuser",
				Role: "",
			},
			valid: false,
			err:   ErrEmptyRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Valid()

			if tt.valid {
				if err != nil {
					t.Errorf("%s: expected no error, but got %s", tt.name, err.Error())
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected %s, got nil", tt.name, tt.err.Error())
				}

				if !errors.Is(err, tt.err) {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err.Error(), err.Error())
				}
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	authenticator := NewJWTAuthenticator("mysecretkey", 24)

	tests := []struct {
		name  string
		token *TokenPayload
		valid bool
		err   error
	}{
		{
			name: "Valid User",
			token: &TokenPayload{
				ID:   1,
				Name: "testuser",
				Role: "user",
			},
			valid: true,
		},
		{
			name: "Invalid User",
			token: &TokenPayload{
				ID:   0,
				Name: "",
				Role: "",
			},
			valid: false,
			err:   ErrZeroID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authenticator.GenerateToken(tt.token)

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

func generateTestToken(secret string, expiryHours int, payload *TokenPayload) string {
	claims := jwt.MapClaims{
		"user": payload,
		"exp":  time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func TestValidateToken(t *testing.T) {
	authenticator := NewJWTAuthenticator("mysecretkey", 24)
	wrongSecretAuthenticator := NewJWTAuthenticator("wrongsecret", 24)

	validPayload := &TokenPayload{
		ID:   1,
		Name: "testuser",
		Role: "admin",
	}

	tests := []struct {
		name           string
		token          string
		expectedClaims jwt.MapClaims
		valid          bool
		err            error
	}{
		{
			name:  "Valid Token",
			token: generateTestToken(authenticator.secret, 1, validPayload),
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
			name:           "Malformed Token",
			token:          "token",
			expectedClaims: nil,
			valid:          false,
			err:            jwt.ErrTokenMalformed,
		},
		{
			name:           "Expired Token",
			token:          generateTestToken(authenticator.secret, -1, validPayload),
			expectedClaims: nil,
			valid:          false,
			err:            jwt.ErrTokenExpired,
		},
		{
			name:           "Wrong Signature Token",
			token:          generateTestToken(wrongSecretAuthenticator.secret, 1, validPayload),
			expectedClaims: nil,
			valid:          false,
			err:            jwt.ErrTokenSignatureInvalid,
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

func TestGetClaimsFromRequest(t *testing.T) {
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
				token := generateTestToken(authenticator.secret, 1, &TokenPayload{
					ID:   1,
					Name: "testuser",
					Role: "admin",
				})
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
			expectedClaims: nil,
			valid:          false,
			err:            jwt.ErrTokenMalformed,
		},
		{
			name: "Request without Token",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				return r
			}(),
			expectedClaims: nil,
			valid:          false,
			err:            http.ErrNoCookie,
		},
		{
			name:           "Nil Request",
			request:        nil,
			expectedClaims: nil,
			valid:          false,
			err:            ErrNilRequest,
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
