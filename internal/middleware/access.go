package middleware

import (
	"forum-project/internal/config"
	"forum-project/internal/models"
	"net/http"
	"strconv"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.AuthorizedUser)
		role := user.Role

		if role != "admin" {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func IsPostAuthor(app *config.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				app.ErrorResponder.BadRequest(rw, "Invalid Post ID", err)
				return
			}

			user := r.Context().Value("user").(*models.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				app.ErrorResponder.NotFound(rw, "Post Not Found", err)
				return
			}

			isAuthorOrAdmin := app.PostService.VerifyPostAuthor(post, user.ID)
			if !isAuthorOrAdmin {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func IsPostAuthorOrAdmin(app *config.AppConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				app.ErrorResponder.BadRequest(rw, "Invalid Post ID", err)
				return
			}

			user := r.Context().Value("user").(*models.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				app.ErrorResponder.NotFound(rw, "Post Not Found", err)
				return
			}

			isAuthorOrAdmin := app.PostService.VerifyPostAuthorOrAdmin(post, user.ID, user.Role)
			if !isAuthorOrAdmin {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
