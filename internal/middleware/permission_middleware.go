package middleware

import (
	"net/http"
	"simple-forum/internal/app"
	"simple-forum/internal/model"
	"strconv"
)

func PermissionMiddleware(app *app.App, permissions ...string) func(http.Handler) http.HandlerFunc {
	requiredPerms := make(map[string]struct{})
	for _, p := range permissions {
		requiredPerms[p] = struct{}{}
	}

	return func(next http.Handler) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(*model.AuthorizedUser)

			stringPostID := r.PathValue("postID")
			id, err := strconv.Atoi(stringPostID)
			if err != nil {
				http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
				app.Logger.Error(err.Error(), "method", r.Method, "status",
					http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
				return
			}

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				http.Error(rw, "Post Not Found", http.StatusNotFound)
				app.Logger.Error(err.Error(), "method", r.Method, "status",
					http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
				return
			}

			if _, ok := requiredPerms["admin"]; ok {
				if user.Role == "admin" {
					next.ServeHTTP(rw, r)
					return
				}
			}

			if _, ok := requiredPerms["author"]; ok {
				if post.AuthorId == user.ID {
					next.ServeHTTP(rw, r)
					return
				}
			}

			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}
	}
}
