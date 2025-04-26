package middleware

import (
	"forum-project/internal/service"
	"net/http"
)

func EnsureAdmin(next http.Handler, u *service.UserService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		result := r.Context().Value("email").(string)

		user, err := u.GetUserByMail(result)
		if err != nil {
			http.Error(rw, "User not found", http.StatusNotFound)
		}
		if user.Role != "admin" {
			http.Error(rw, "User is not admin", http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
