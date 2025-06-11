package handler

import (
	"errors"
	"forum-project/internal/app"
	"forum-project/internal/model"
	"forum-project/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	app *app.App
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{app: app}
}

func (u *UserHandler) GetRegister(rw http.ResponseWriter, r *http.Request) {
	err := u.app.Templates.Render(rw, r, "register.page", new(model.Page))
	if err != nil {
		u.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) PostRegister(rw http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")

	err := u.app.Authenticator.Register(username, email, password1, password2)
	if err != nil {
		var errorMsg string
		switch {
		case errors.Is(err, service.ErrUserEmailAlreadyExists):
			errorMsg = "Email already exists"
		case errors.Is(err, service.ErrUserNameAlreadyExists):
			errorMsg = "Username already exists"
		default:
			errorMsg = "Failed to register"
		}
		page := &model.Page{
			Error: errorMsg,
		}
		err := u.app.Templates.Render(rw, r, "register.page", page)
		if err != nil {
			u.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (u *UserHandler) GetLogin(rw http.ResponseWriter, r *http.Request) {
	err := u.app.Templates.Render(rw, r, "login.page", new(model.Page))
	if err != nil {
		u.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) PostLogin(rw http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	token, err := u.app.Authenticator.Login(email, password)
	if err != nil {
		var errorMsg string
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			errorMsg = "User not found"
		case errors.Is(err, service.ErrWrongPassword):
			errorMsg = "Wrong password"
		default:
			errorMsg = "Failed to login"
		}
		page := &model.Page{
			Error: errorMsg,
		}
		err := u.app.Templates.Render(rw, r, "login.page", page)
		if err != nil {
			u.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
			return
		}
		return
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

func (u *UserHandler) handleError(rw http.ResponseWriter, msg string, err error, code int) {
	http.Error(rw, msg, code)
	u.app.Logger.Error(msg, "error", err.Error())
}
