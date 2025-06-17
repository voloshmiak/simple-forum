package middleware

import (
	"forum-project/internal/app"
	"forum-project/internal/model"
	"forum-project/internal/service"
	"net/http"
	"strconv"
)

func RBACMiddleware(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		})
	}
}

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

func IsPostAuthor(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
				app.Logger.Error(err.Error(), "method", r.Method, "status",
					http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
				return
			}

			user := r.Context().Value("user").(*model.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				http.Error(rw, "Post Not Found", http.StatusNotFound)
				app.Logger.Error(err.Error(), "method", r.Method, "status",
					http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
				return
			}

			isAuthor := service.VerifyPostAuthor(post, user.ID)
			if !isAuthor {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func IsPostAuthorOrAdmin(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
				app.Logger.Error(err.Error(), "method", r.Method, "status",
					http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
				return
			}

			user := r.Context().Value("user").(*model.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				http.Error(rw, "Post Not Found", http.StatusNotFound)
				app.Logger.Error(err.Error(), "method", r.Method, "status",
					http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
				return
			}

			isAuthorOrAdmin := service.VerifyPostAuthorOrAdmin(post, user.ID, user.Role)
			if !isAuthorOrAdmin {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
