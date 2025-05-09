package handlers

import (
	"fmt"
	"forum-project/internal/auth"
	"forum-project/internal/config"
	"forum-project/internal/models"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type PostHandler struct {
	app *config.AppConfig
}

func NewPostHandler(app *config.AppConfig) *PostHandler {
	return &PostHandler{app: app}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Errors.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Post Not Found", err)
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

		isAuthor := p.app.PostService.VerifyPostAuthor(post, userIDInt)
		if isAuthor {
			viedData.IsAuthor = true
		}
	}

	data := make(map[string]any)
	data["post"] = post

	viedData.Data = data

	err = p.app.Templates.Render(rw, r, "post.page", viedData)
	if err != nil {
		p.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) GetCreatePost(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		p.app.Errors.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	topic, err := p.app.TopicService.GetTopicByID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Topic Not Found", err)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = p.app.Templates.Render(rw, r, "create-post.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		p.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) PostCreatePost(rw http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	topicID := r.PostFormValue("topic_id")
	topicIDInt, err := strconv.Atoi(topicID)
	if err != nil {
		p.app.Errors.BadRequest(rw, "Invalid Topic ID", err)
		return
	}

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)
	userName := user.(map[string]interface{})["username"].(string)

	err = p.app.PostService.CreatePost(title, content, topicIDInt, userID, userName)
	if err != nil {
		p.app.Errors.InternalServer(rw, "Unable to create post", err)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%d", topicIDInt)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Errors.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Post Not Found", err)
		return
	}

	isAuthor := p.app.PostService.VerifyPostAuthor(post, userID)
	if !isAuthor {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	data := make(map[string]any)
	data["post"] = post

	err = p.app.Templates.Render(rw, r, "edit-post.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		p.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (p *PostHandler) PostEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Errors.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Post Not Found", err)
		return
	}

	isAuthor := p.app.PostService.VerifyPostAuthor(post, userID)
	if !isAuthor {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Topic Not Found", err)
		return
	}

	post.Title = title
	post.Content = content

	err = p.app.PostService.EditPost(post)
	if err != nil {
		p.app.Errors.InternalServer(rw, "Unable to edit post", err)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetDeletePost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.app.Errors.BadRequest(rw, "Invalid Post ID", err)
		return
	}

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Topic Not Found", err)
		return
	}

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)
	userRole := user.(map[string]interface{})["role"].(string)

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.app.Errors.NotFound(rw, "Post Not Found", err)
		return
	}

	isAuthorOrAdmin := p.app.PostService.VerifyPostAuthorOrAdmin(post, userID, userRole)
	if !isAuthorOrAdmin {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	err = p.app.PostService.DeletePost(id)
	if err != nil {
		p.app.Errors.InternalServer(rw, "Unable to delete post", err)
		return
	}

	url := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, url, http.StatusFound)
}
