package middleware

import (
	"forum-project/internal/application"
	"forum-project/internal/model"
	"net/http"
	"strconv"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*model.AuthorizedUser)
		role := user.Role

		if role != "admin" {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func IsPostAuthor(app *application.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
				return
			}

			user := r.Context().Value("user").(*model.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				http.Error(rw, "Post Not Found", http.StatusNotFound)
				return
			}

			isAuthor := app.PostService.VerifyPostAuthor(post, user.ID)
			if !isAuthor {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func IsPostAuthorOrAdmin(app *application.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
				return
			}

			user := r.Context().Value("user").(*model.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				http.Error(rw, "Post Not Found", http.StatusNotFound)
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
