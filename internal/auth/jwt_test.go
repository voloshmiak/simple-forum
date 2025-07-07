package auth

import (
	"net/http"
	"net/http/httptest"
	"simple-forum/internal/model"
	"testing"
	"time"
)

var (
	testUser = &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@email.com",
		PasswordHash: "testpassword",
		CreatedAt:    time.Now(),
		Role:         "user",
	}
	authenticator = NewJWTAuthenticator("mysecretkey", 24)
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name  string
		user  *model.User
		valid bool
		err   string
	}{
		{
			name:  "Valid User",
			user:  testUser,
			valid: true,
			err:   "",
		},
		{
			name:  "Nil User",
			user:  nil,
			valid: false,
			err:   "user cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			token, err := authenticator.GenerateToken(tt.user)

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
				if token != "" {
					t.Errorf("%s: expected no token, but got %s", tt.name, token)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	validToken, _ := authenticator.GenerateToken(testUser)
	expiredToken, _ := NewJWTAuthenticator("mysecretkey", -1).GenerateToken(testUser)
	wrongSecretToken, _ := NewJWTAuthenticator("wrongsecret", 24).GenerateToken(testUser)

	tests := []struct {
		name  string
		token string
		valid bool
		err   string
	}{
		{
			name:  "Valid Token",
			token: validToken,
			valid: true,
			err:   "",
		},
		{
			name:  "Malformed Token",
			token: "token",
			valid: false,
			err:   "token is malformed: token contains an invalid number of segments",
		},
		{
			name:  "Expired Token",
			token: expiredToken,
			valid: false,
			err:   "token has invalid claims: token is expired",
		},
		{
			name:  "Token with wrong signature",
			token: wrongSecretToken,
			valid: false,
			err:   "token signature is invalid: signature is invalid",
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
			} else {
				if err == nil {
					t.Errorf("%s: expected %s, got nil", tt.name, tt.err)
				}
				if claims != nil {
					t.Errorf("%s: expected no claims, but got %s", tt.name, claims)
				}
			}
		})
	}
}

func TestGetClaimsFromRequest(t *testing.T) {
	token, _ := authenticator.GenerateToken(testUser)

	validRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validRequest.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	invalidRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidRequest.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    "invalidtoken",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	tests := []struct {
		name    string
		request *http.Request
		valid   bool
		err     string
	}{
		{
			name:    "Valid Request with Token",
			request: validRequest,
			valid:   true,
			err:     "",
		},
		{
			name:    "Nil request",
			request: nil,
			valid:   false,
			err:     "http: named cookie not present",
		},
		{
			name:    "Request without Token",
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			valid:   false,
			err:     "http: named cookie not present",
		},
		{
			name:    "Request with Invalid Token",
			request: invalidRequest,
			valid:   false,
			err:     "token is malformed: token contains an invalid number of segments",
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
			} else {
				if err == nil {
					t.Errorf("%s: expected an error, got nil", tt.name)
				}
				if claims != nil {
					t.Errorf("%s: expected no claims, but got %s", tt.name, claims)
				}
			}
		})
	}
}
