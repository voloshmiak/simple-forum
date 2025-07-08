package template

import (
	"net/http"
	"net/http/httptest"
	"simple-forum/internal/auth"
	"simple-forum/internal/model"
	"strings"
	"testing"
	"time"
)

func TestRender(t *testing.T) {
	authenticator := auth.NewJWTAuthenticator("mysecretkey", 24)
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

	adminToken, err := authenticator.GenerateToken(1, "testuser", "admin")
	if err != nil {
		t.Fatal(err)
	}
	adminRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	adminRequest.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    adminToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	tests := []struct {
		name    string
		writer  http.ResponseWriter
		request *http.Request
		tmpl    string
		valid   bool
		page    *model.Page
		err     string
	}{
		{
			name:    "Valid Template",
			writer:  httptest.NewRecorder(),
			request: validRequest,
			tmpl:    "home.page",
			page:    new(model.Page),
			valid:   true,
		},
		{
			name:    "Valid Admin Template",
			writer:  httptest.NewRecorder(),
			request: adminRequest,
			tmpl:    "home.page",
			page:    new(model.Page),
			valid:   true,
		},
		{
			name:    "Invalid Template",
			writer:  httptest.NewRecorder(),
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			tmpl:    "invalid.template",
			page:    new(model.Page),
			valid:   false,
			err:     "invalid.template.gohtml not found",
		},
		{
			name:    "Nil ResponseWriter",
			writer:  nil,
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			tmpl:    "home.page",
			page:    new(model.Page),
			valid:   false,
			err:     "ResponseWriter is nil",
		},
		{
			name:    "Nil Request",
			writer:  httptest.NewRecorder(),
			request: nil,
			tmpl:    "home.page",
			page:    new(model.Page),
			valid:   false,
			err:     "Request is nil",
		},
		{
			name:    "Nil Page",
			writer:  httptest.NewRecorder(),
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			tmpl:    "home.page",
			page:    nil,
			valid:   true,
		},
	}

	templates := NewTemplates("development", "../../web/templates", authenticator)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				rr := tt.writer.(*httptest.ResponseRecorder)

				err := templates.Render(rr, tt.request, tt.tmpl, tt.page)
				if err != nil {
					t.Errorf("%s: expected no error, got %s", tt.name, err.Error())
				}

				text := `<h1 class="cover-heading">Explore the Forum.</h1>`

				if !strings.Contains(rr.Body.String(), text) {
					t.Errorf("%s: expected %s in the response, got %s", tt.name, text, rr.Body.String())
				}
			} else {
				err := templates.Render(tt.writer, tt.request, tt.tmpl, tt.page)
				if err == nil {
					t.Errorf("%s: expected an error, got nil", tt.name)
				}
				if tt.err != err.Error() {
					t.Errorf("%s: expected %s, got %s", tt.name, tt.err, err.Error())
				}
			}
		})
	}
}
