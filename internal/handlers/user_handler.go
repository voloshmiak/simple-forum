package handlers

import (
	"errors"
	"forum-project/internal/config"
	"forum-project/internal/models"
	"forum-project/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	app *config.AppConfig
}

func NewUserHandler(app *config.AppConfig) *UserHandler {
	return &UserHandler{app: app}
}

func (u *UserHandler) GetRegister(rw http.ResponseWriter, r *http.Request) {
	err := u.app.Templates.Render(rw, r, "register.page", &models.Page{})
	if err != nil {
		u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
	}
}

func (u *UserHandler) PostRegister(rw http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")

	err := u.app.UserService.Register(username, email, password1, password2)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserEmailAlreadyExists):
			page := &models.Page{
				Error: "Email already exists",
			}
			err := u.app.Templates.Render(rw, r, "register.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		case errors.Is(err, service.ErrUserNameAlreadyExists):
			page := &models.Page{
				Error: "Username already exists",
			}
			err := u.app.Templates.Render(rw, r, "register.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		case errors.Is(err, service.ErrMissmatchPassword):
			page := &models.Page{
				Error: "Passwords do not match",
			}
			err := u.app.Templates.Render(rw, r, "register.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		default:
			page := &models.Page{
				Error: "Failed to register",
			}
			err := u.app.Templates.Render(rw, r, "register.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		}
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (u *UserHandler) GetLogin(rw http.ResponseWriter, r *http.Request) {
	err := u.app.Templates.Render(rw, r, "login.page", new(models.Page))
	if err != nil {
		u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
	}
}

func (u *UserHandler) PostLogin(rw http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	token, err := u.app.UserService.Authenticate(email, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			page := &models.Page{
				Error: "User not found",
			}
			err := u.app.Templates.Render(rw, r, "login.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		case errors.Is(err, service.ErrWrongPassword):
			page := &models.Page{
				Error: "Wrong password",
			}
			err := u.app.Templates.Render(rw, r, "login.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		default:
			page := &models.Page{
				Error: "Failed to login",
			}
			err := u.app.Templates.Render(rw, r, "login.page", page)
			if err != nil {
				u.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
			}
			return
		}
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

	http.Redirect(rw, r, "/topics", http.StatusFound)

}

func (u *UserHandler) GetLogout(rw http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(-time.Hour * 24),
	}

	http.SetCookie(rw, cookie)

	http.Redirect(rw, r, "/home", http.StatusFound)
}
