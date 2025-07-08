package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"simple-forum/internal/auth"
	"simple-forum/internal/model"
	"simple-forum/internal/service"
	"simple-forum/internal/template"
	"time"
)

type UserService interface {
	Login(email, password string) (*model.User, error)
	Register(username, email, password1, password2 string) error
}

type UserHandler struct {
	l  *slog.Logger
	a  *auth.JWTAuthenticator
	t  *template.Templates
	us UserService
}

func NewUserHandler(l *slog.Logger, a *auth.JWTAuthenticator, t *template.Templates, us UserService) *UserHandler {
	return &UserHandler{l: l, t: t, a: a, us: us}
}

func (u *UserHandler) GetRegister(rw http.ResponseWriter, r *http.Request) {
	err := u.t.Render(rw, r, "register.page", new(model.Page))
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		u.l.Error(msg, "error", err.Error())
		return
	}
}

func (u *UserHandler) PostRegister(rw http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")

	err := u.us.Register(username, email, password1, password2)
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
		err = u.t.Render(rw, r, "register.page", page)
		if err != nil {
			msg := "Unable to render template"
			http.Error(rw, msg, http.StatusInternalServerError)
			u.l.Error(msg, "error", err.Error())
			return
		}
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (u *UserHandler) GetLogin(rw http.ResponseWriter, r *http.Request) {
	err := u.t.Render(rw, r, "login.page", new(model.Page))
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		u.l.Error(msg, "error", err.Error())
		return
	}
}

func (u *UserHandler) PostLogin(rw http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	user, err := u.us.Login(email, password)
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
		err = u.t.Render(rw, r, "login.page", page)
		if err != nil {
			msg := "Unable to render template"
			http.Error(rw, msg, http.StatusInternalServerError)
			u.l.Error(msg, "error", err.Error())
			return
		}
		return
	}

	token, err := u.a.GenerateToken(user.ID, user.Name, user.Role)
	if err != nil {
		msg := "Failed to generate token"
		http.Error(rw, msg, http.StatusInternalServerError)
		u.l.Error(msg, "error", err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
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
