package handlers

import (
	"fmt"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

type AuthHandler struct {
	logger      *slog.Logger
	templates   *template.Manager
	userService *service.UserService
}

func NewAuthHandler(logger *slog.Logger, templates *template.Manager, userService *service.UserService) *AuthHandler {
	return &AuthHandler{logger: logger, templates: templates, userService: userService}
}

func (a *AuthHandler) GetRegister(rw http.ResponseWriter, r *http.Request) {
	err := a.templates.Render(rw, "register.page", nil)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (a *AuthHandler) PostRegister(rw http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")
	_, err := a.userService.Register(username, email, password1, password2)
	if err != nil {
		rw.Write([]byte(err.Error()))
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (a *AuthHandler) GetLogin(rw http.ResponseWriter, r *http.Request) {
	err := a.templates.Render(rw, "login.page", nil)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (a *AuthHandler) PostLogin(rw http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	user, err := a.userService.Authenticate(email, password)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to login: %s", err), http.StatusInternalServerError)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte("secret-key"))

	if err != nil {
		rw.Write([]byte(err.Error()))
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(rw, cookie)

	a.logger.Info(fmt.Sprintf("User authenticated: %s", user))

	http.Redirect(rw, r, "/topics", http.StatusFound)

}

func (a *AuthHandler) GetLogout(rw http.ResponseWriter, r *http.Request) {
	err := a.templates.Render(rw, "logout.page", nil)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (a *AuthHandler) PostLogout(rw http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(-time.Hour * 24),
	}

	http.SetCookie(rw, cookie)

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
