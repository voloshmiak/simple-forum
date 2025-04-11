package handlers

import (
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
	p.logger.Info("Handle GET Posts")

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.Error("Unable to convert id to integer", err)
		http.Error(rw, "Unable to convert id to integer", http.StatusBadRequest)
		return
	}

	posts, err := models.GetTopicPosts(id)
	if err != nil {
		p.logger.Error("Unable to get posts", err)
		http.Error(rw, "Unable to get posts", http.StatusInternalServerError)
		return
	}

	err = p.renderer.RenderTemplate(rw, "posts.page", posts)
	if err != nil {
		p.logger.Error("Unable to render template", err)
	}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {
	p.logger.Info("Handle GET Post")

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		http.Redirect(rw, r, "/posts/", http.StatusFound)
		return
	}

	post, err := models.FindPost(id)
	if err != nil {
		p.logger.Error("Unable to find post", err)
		http.Error(rw, "Unable to find post", http.StatusNotFound)
		return
	}

	err = p.renderer.RenderTemplate(rw, "post.page", post)
	if err != nil {
		p.logger.Error("Unable to render template", err)
	}
}
