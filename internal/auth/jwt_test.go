package auth

import (
	"net/http"
	"net/http/httptest"
	"simple-forum/internal/model"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		user        *model.User
		expectToken bool
		err         string
	}{
		{
			name: "Valid User",
			user: &model.User{
				ID:           1,
				Username:     "testuser",
				Email:        "test@email.com",
				PasswordHash: "testpassword",
				CreatedAt:    time.Now(),
				Role:         "user",
			},
			expectToken: true,
			err:         "",
		},
		{
			name:        "Nil User",
			user:        nil,
			expectToken: false,
			err:         "user cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := "mysecretkey"
			expiryHours := 24

			authenticator := NewJWTAuthenticator(secret, expiryHours)

			token, err := authenticator.GenerateToken(tt.user)

			if tt.expectToken {
				if err != nil {
					t.Errorf("%s: no error expected, but got %s", tt.name, err.Error())
				}
				if token == "" {
					t.Errorf("%s: expected a token, got an empty string", tt.name)
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected an error, got nil", tt.name)
				}
				if err.Error() != tt.err {
					t.Errorf("%s: expected %s, but got %s", tt.name, tt.err, err.Error())
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	testUser := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@email.com",
		PasswordHash: "testpassword",
		CreatedAt:    time.Now(),
		Role:         "user",
	}

	secret := "mysecretkey"
	expiryHours := 24

	authenticator := NewJWTAuthenticator(secret, expiryHours)

	token, _ := authenticator.GenerateToken(testUser)

	tests := []struct {
		name        string
		token       string
		expectValid bool
		err         string
	}{
		{
			name:        "Valid Token",
			token:       token,
			expectValid: true,
			err:         "",
		},
		{
			name:        "Invalid Token",
			token:       "invalidtoken",
			expectValid: false,
			err:         "token is malformed: token contains an invalid number of segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authenticator.validateToken(tt.token)

			if tt.expectValid {
				if err != nil {
					t.Errorf("%s: no error expected, but got %s", tt.name, err.Error())
				}
				if claims == nil {
					t.Errorf("%s: expected claims, got nil", tt.name)
				} else {
					user := claims["user"].(map[string]interface{})
					userIDFloat := user["id"].(float64)
					userIDInt := int(userIDFloat)
					if userIDInt != testUser.ID {
						t.Errorf("%s: expected user ID %d, got %v", tt.name, testUser.ID, userIDInt)
					}
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected an error, got nil", tt.name)
				}
				if err.Error() != tt.err {
					t.Errorf("%s: expected %s, but got %s", tt.name, tt.err, err.Error())
				}
			}
		})
	}
}

func TestGetClaimsFromRequest(t *testing.T) {
	testUser := &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@email.com",
		PasswordHash: "testpassword",
		CreatedAt:    time.Now(),
		Role:         "user",
	}

	secret := "mysecretkey"
	expiryHours := 24

	authenticator := NewJWTAuthenticator(secret, expiryHours)

	token, _ := authenticator.GenerateToken(testUser)

	tests := []struct {
		name        string
		token       string
		needCookie  bool
		expectClaim bool
		err         string
	}{
		{
			name:        "Valid Request with Token",
			token:       token,
			needCookie:  true,
			expectClaim: true,
			err:         "",
		},
		{
			name:        "Valid Request with Invalid Token",
			token:       "invalidtoken",
			needCookie:  true,
			expectClaim: false,
			err:         "token is malformed: token contains an invalid number of segments",
		},
		{
			name:        "Request without Token",
			token:       "",
			needCookie:  false,
			expectClaim: false,
			err:         "http: named cookie not present",
		},
	}

	for _, tt := range tests {
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		if tt.needCookie {
			cookie := &http.Cookie{
				Name:     "token",
				Value:    tt.token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  time.Now().Add(time.Hour * 24),
			}

			request.AddCookie(cookie)
		}

		t.Run(tt.name, func(t *testing.T) {
			claims, err := authenticator.GetClaimsFromRequest(request)

			if tt.expectClaim {
				if err != nil {
					t.Errorf("%s: no error expected, but got %s", tt.name, err.Error())
				}
				if claims == nil {
					t.Errorf("%s: expected claims, got nil", tt.name)
				} else {
					user := claims["user"].(map[string]interface{})
					userIDFloat := user["id"].(float64)
					userIDInt := int(userIDFloat)
					if userIDInt != testUser.ID {
						t.Errorf("%s: expected user ID %d, got %v", tt.name, testUser.ID, userIDInt)
					}
				}
			} else {
				if err == nil {
					t.Errorf("%s: expected an error, got nil", tt.name)
				}
				if err.Error() != tt.err {
					t.Errorf("%s: expected %s, but got %s", tt.name, tt.err, err.Error())
				}
			}
		})
	}
}
