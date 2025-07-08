package middleware

import (
	"log/slog"
	"net/http"
	"simple-forum/internal/handler"
	"strconv"
)

func PermissionMiddleware(l *slog.Logger, ps handler.PostService, permissions ...string) func(http.Handler) http.HandlerFunc {
	requiredPerms := make(map[string]struct{})
	for _, p := range permissions {
		requiredPerms[p] = struct{}{}
	}

	return func(next http.Handler) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			userValue := r.Context().Value("user")
			if userValue == nil {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, ok := userValue.(map[string]interface{})
			if !ok {
				http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
				l.Error("Invalid user type in context", "method", r.Method, "path", r.URL.Path)
				return
			}

			userIDFloat, ok := user["id"].(float64)
			if !ok {
				http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
				l.Error("Invalid user ID type", "method", r.Method, "path", r.URL.Path)
				return
			}

			userRole, ok := user["role"].(string)
			if !ok {
				http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
				l.Error("Invalid user role type", "method", r.Method, "path", r.URL.Path)
				return
			}

			userID := int(userIDFloat)

			if _, ok := requiredPerms["admin"]; ok {
				if userRole == "admin" {
					next.ServeHTTP(rw, r)
					return
				}
			}

			if _, ok := requiredPerms["author"]; ok {
				stringPostID := r.PathValue("postID")
				id, err := strconv.Atoi(stringPostID)
				if err != nil {
					http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
					l.Error(err.Error(), "method", r.Method, "status",
						http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
					return
				}

				post, err := ps.GetPostByID(id)
				if err != nil {
					http.Error(rw, "Post Not Found", http.StatusNotFound)
					l.Error(err.Error(), "method", r.Method, "status",
						http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
					return
				}

				if post.AuthorId == userID {
					next.ServeHTTP(rw, r)
					return
				}
			}

			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}
	}
}
