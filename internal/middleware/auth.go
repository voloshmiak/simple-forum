package middleware

import (
	"context"
	"forum-project/internal/auth"
	"forum-project/internal/config"
	"forum-project/internal/models"
	"net/http"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		claims, err := auth.GetClaimsFromRequest(r)
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user := claims["user"].(map[string]interface{})

		userIDFloat := user["id"].(float64)
		userIDInt := int(userIDFloat)
		username := user["username"].(string)
		email := user["email"].(string)
		role := user["role"].(string)

		authorizedUser := &models.AuthorizedUser{
			ID:       userIDInt,
			Username: username,
			Email:    email,
			Role:     role,
		}

		ctx := context.WithValue(r.Context(), "user", authorizedUser)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

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
				app.Errors.BadRequest(rw, "Invalid Post ID", err)
				return
			}

			user := r.Context().Value("user").(*models.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				app.Errors.NotFound(rw, "Post Not Found", err)
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
				app.Errors.BadRequest(rw, "Invalid Post ID", err)
				return
			}

			user := r.Context().Value("user").(*models.AuthorizedUser)

			post, err := app.PostService.GetPostByID(id)
			if err != nil {
				app.Errors.NotFound(rw, "Post Not Found", err)
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
