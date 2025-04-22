package handlers

import (
	"fmt"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"strconv"
)

type PostHandler struct {
	logger      *slog.Logger
	templates   *template.Manager
	postService *service.PostService
}

func NewPostHandler(logger *slog.Logger, renderer *template.Manager, postService *service.PostService) *PostHandler {
	return &PostHandler{logger, renderer, postService}
}

func (p *PostHandler) GetPostsByTopicID(rw http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.Error("Unable to convert id to integer")
		http.Error(rw, "Unable to convert id to integer", http.StatusBadRequest)
		return
	}

	posts, err := p.postService.GetPostsByTopicID(id)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Unable to get posts: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to get posts: %s", err), http.StatusInternalServerError)
		return
	}

	err = p.templates.Render(rw, "posts.page", posts)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (p *PostHandler) GetPostByID(rw http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		http.Redirect(rw, r, "/posts/", http.StatusFound)
		return
	}

	post, err := p.postService.GetPostByID(id)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Post not found: %s", err))
		http.Error(rw, fmt.Sprintf("Post not found: %s", err), http.StatusNotFound)
		return
	}

	err = p.templates.Render(rw, "post.page", post)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (p *PostHandler) CreatePost(rw http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) UpdatePost(rw http.ResponseWriter, r *http.Request) {}

func (p *PostHandler) DeletePost(rw http.ResponseWriter, r *http.Request) {}
