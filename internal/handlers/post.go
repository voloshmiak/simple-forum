package handlers

import (
	"fmt"
	"forum-project/internal/auth"
	"forum-project/internal/models"
	"forum-project/internal/mylogger"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type PostHandler struct {
	logger       *mylogger.WrappedLogger
	templates    *template.Manager
	postService  *service.PostService
	topicService *service.TopicService
}

func NewPostHandler(logger *mylogger.WrappedLogger, renderer *template.Manager, postService *service.PostService, topicService *service.TopicService) *PostHandler {
	return &PostHandler{logger, renderer, postService, topicService}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("postID")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.BadRequestError(rw, "Invalid Post ID", err)
		return
	}

	post, err := p.postService.GetPostByID(id)
	if err != nil {
		p.logger.NotFoundError(rw, "Post Not Found", err)
		return
	}

	viedData := &models.ViewData{}
	viedData.IsAuthor = false

	token, err := auth.ValidateTokenFromRequest(r)
	if err == nil {
		claims := token.Claims.(jwt.MapClaims)

		userClaim := claims["user"].(map[string]interface{})
		userIDClaim := userClaim["id"].(float64)
		userIDInt := int(userIDClaim)

		if post.AuthorId == userIDInt {
			viedData.IsAuthor = true
		}
	}

	data := make(map[string]any)
	data["post"] = post

	viedData.Data = data

	err = p.templates.Render(rw, r, "post.page", viedData)
	if err != nil {
		p.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) GetCreatePost(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.BadRequestError(rw, "Invalid Post ID", err)
		return
	}

	topic, err := p.topicService.GetTopicByID(topicID)
	if err != nil {
		p.logger.NotFoundError(rw, "Topic Not Found", err)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = p.templates.Render(rw, r, "create-post.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		p.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) PostCreatePost(rw http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	topicID := r.PostFormValue("topic_id")
	topicIDInt, err := strconv.Atoi(topicID)
	if err != nil {
		p.logger.BadRequestError(rw, "Invalid Topic ID", err)
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	user := claims["user"].(map[string]interface{})
	authorIDfloat := user["id"].(float64)
	authorIDInt := int(authorIDfloat)

	_, err = p.postService.CreatePost(title, content, topicIDInt, authorIDInt)
	if err != nil {
		p.logger.ServerInternalError(rw, "Unable to create post", err)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%d", topicIDInt)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetEditPost(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	postID, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.BadRequestError(rw, "Invalid Post ID", err)
		return
	}

	post, err := p.postService.GetPostByID(postID)
	if err != nil {
		p.logger.NotFoundError(rw, "Post Not Found", err)
		return
	}

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)

	if userID != post.AuthorId {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	data := make(map[string]any)
	data["post"] = post

	err = p.templates.Render(rw, r, "edit-post.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		p.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) PostEditPost(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	postID, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.BadRequestError(rw, "Invalid Post ID", err)
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	post, err := p.postService.GetPostByID(postID)
	if err != nil {
		p.logger.NotFoundError(rw, "Post Not Found", err)
		return
	}

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)

	if userID != post.AuthorId {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	topic, err := p.topicService.GetTopicByPostID(postID)
	if err != nil {
		p.logger.NotFoundError(rw, "Topic Not Found", err)
		return
	}

	post.Title = title
	post.Content = content

	_, err = p.postService.EditPost(post)
	if err != nil {
		p.logger.ServerInternalError(rw, "Unable to edit post", err)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetDeletePost(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		p.logger.BadRequestError(rw, "Invalid Post ID", err)
		return
	}

	topic, err := p.topicService.GetTopicByPostID(id)
	if err != nil {
		p.logger.NotFoundError(rw, "Topic Not Found", err)
		return
	}

	post, err := p.postService.GetPostByID(id)
	if err != nil {
		p.logger.NotFoundError(rw, "Post Not Found", err)
		return
	}

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)

	userRole := user.(map[string]interface{})["role"].(string)

	if userID != post.AuthorId && userRole != "admin" {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	err = p.postService.DeletePost(id)
	if err != nil {
		p.logger.ServerInternalError(rw, "Unable to delete post", err)
		return
	}

	url := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, url, http.StatusFound)
}
