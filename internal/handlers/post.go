package handlers

import (
	"fmt"
	"forum-project/internal/models"
	"forum-project/internal/render"
	"log/slog"
	"net/http"
	"strconv"
)

type PostHandler struct {
	logger   *slog.Logger
	renderer *render.Renderer
}

func NewPostHandler(logger *slog.Logger, renderer *render.Renderer) *PostHandler {
	return &PostHandler{logger, renderer}
}

func (p *PostHandler) GetPosts(rw http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.Error("Unable to convert id to integer")
		http.Error(rw, "Unable to convert id to integer", http.StatusBadRequest)
		return
	}

	posts, err := models.GetTopicPosts(id)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Unable to get posts: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to get posts: %s", err), http.StatusInternalServerError)
		return
	}

	err = p.renderer.RenderTemplate(rw, "posts.page", posts)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Unable to render template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to render template: %s", err), http.StatusInternalServerError)
	}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		http.Redirect(rw, r, "/posts/", http.StatusFound)
		return
	}

	post, err := models.FindPost(id)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Post not found: %s", err))
		http.Error(rw, fmt.Sprintf("Post not found: %s", err), http.StatusNotFound)
		return
	}

	err = p.renderer.RenderTemplate(rw, "post.page", post)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Unable to render template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to render template: %s", err), http.StatusInternalServerError)
	}
}

func (p *PostHandler) CreatePost(rw http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) UpdatePost(rw http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) DeletePost(rw http.ResponseWriter, r *http.Request) {}
