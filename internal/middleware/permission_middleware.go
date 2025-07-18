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

			msg := "Failed to get user"

			user, ok := userValue.(map[string]interface{})
			if !ok {
				http.Error(rw, msg, http.StatusInternalServerError)
				l.Error(msg, "method", r.Method, "path", r.URL.Path)
				return
			}

			userIDFloat, ok := user["id"].(float64)
			if !ok {
				http.Error(rw, msg, http.StatusInternalServerError)
				l.Error(msg, "method", r.Method, "path", r.URL.Path)
				return
			}

			userRole, ok := user["role"].(string)
			if !ok {
				http.Error(rw, msg, http.StatusInternalServerError)
				l.Error(msg, "method", r.Method, "path", r.URL.Path)
				return
			}

			userID := int(userIDFloat)

			if _, ok = requiredPerms["admin"]; ok {
				if userRole == "admin" {
					next.ServeHTTP(rw, r)
					return
				}
			}

			if _, ok = requiredPerms["author"]; ok {
				stringPostID := r.PathValue("postID")
				id, err := strconv.Atoi(stringPostID)
				if err != nil {
					http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
					return
				}

				post, err := ps.GetPostByID(id)
				if err != nil {
					http.Error(rw, "Post Not Found", http.StatusNotFound)
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
