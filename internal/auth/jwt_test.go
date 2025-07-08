package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	authenticator := NewJWTAuthenticator("mysecretkey", 24)

	tests := []struct {
		name     string
		id       int
		username string
		role     string
		valid    bool
		err      string
	}{
		{
			name:     "Valid User",
			id:       1,
			username: "testuser",
			role:     "user",
			valid:    true,
		},
		{
			name:     "Zero User ID",
			id:       0,
			username: "testuser",
			role:     "user",
			valid:    false,
			err:      "id cannot be 0",
		},
		{
			name:     "Empty User name",
			id:       1,
			username: "",
			role:     "user",
			valid:    false,
			err:      "username cannot be empty",
		},
		{
			name:     "Empty User role",
			id:       1,
			username: "testuser",
			role:     "",
			valid:    false,
			err:      "role cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			token, err := authenticator.GenerateToken(tt.id, tt.username, tt.role)

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
				if tt.err != err.Error() {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err, err.Error())
				}
				if token != "" {
					t.Errorf("%s: expected no token, but got %s", tt.name, token)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	authenticator := NewJWTAuthenticator("mysecretkey", 24)
	validToken, err := authenticator.GenerateToken(1, "testuser", "user")
	if err != nil {
		t.Fatal(err)
	}

	expiredToken, err := NewJWTAuthenticator("mysecretkey", -1).GenerateToken(1, "testuser", "user")
	if err != nil {
		t.Fatal(err)
	}

	wrongSecretToken, err := NewJWTAuthenticator("wrongsecret", 24).GenerateToken(1, "testuser", "user")
	if err != nil {
		t.Fatal(err)
	}

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
				if tt.err != err.Error() {
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
	token, err := authenticator.GenerateToken(1, "testuser", "user")
	if err != nil {
		t.Fatal(err)
	}

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
		},
		{
			name:    "Nil request",
			request: nil,
			valid:   false,
			err:     "request cannot be nil",
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
				if tt.err != err.Error() {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err, err.Error())
				}
				if claims != nil {
					t.Errorf("%s: expected no claims, but got %s", tt.name, claims)
				}
			}
		})
	}
}
