package handlers

import (
	"fmt"
	"forum-project/internal/auth"
	"forum-project/internal/models"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"time"
)

type UserHandler struct {
	logger      *slog.Logger
	templates   *template.Manager
	userService *service.UserService
}

func NewUserHandler(logger *slog.Logger, templates *template.Manager, userService *service.UserService) *UserHandler {
	return &UserHandler{logger: logger, templates: templates, userService: userService}
}

func (u *UserHandler) GetRegister(rw http.ResponseWriter, r *http.Request) {
	err := u.templates.Render(rw, r, "register.page", &models.ViewData{})
	if err != nil {
		u.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (u *UserHandler) PostRegister(rw http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")
	_, err := u.userService.Register(username, email, password1, password2)
	if err != nil {
		rw.Write([]byte(err.Error()))
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (u *UserHandler) GetLogin(rw http.ResponseWriter, r *http.Request) {
	err := u.templates.Render(rw, r, "login.page", &models.ViewData{})
	if err != nil {
		u.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (u *UserHandler) PostLogin(rw http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	user, err := u.userService.Authenticate(email, password)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to login: %s", err), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user)

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

	u.logger.Info("User authenticated", "user", user)

	http.Redirect(rw, r, "/topics", http.StatusFound)

}

func (a *UserHandler) GetLogout(rw http.ResponseWriter, r *http.Request) {
	err := a.templates.Render(rw, r, "logout.page", &models.ViewData{})
	if err != nil {
		a.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (a *UserHandler) PostLogout(rw http.ResponseWriter, r *http.Request) {
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
