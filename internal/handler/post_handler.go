package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"simple-forum/internal/auth"
	"simple-forum/internal/model"
	"simple-forum/internal/template"
	"strconv"
)

type PostService interface {
	GetPostByID(postID int) (*model.Post, error)
	GetPostsByTopicID(topicID int) ([]*model.Post, error)
	CreatePost(title, content string, topicID, authorID int, authorName string) error
	EditPost(title, content string, postID int) error
	DeletePost(postID int) error
}

type PostHandler struct {
	l  *slog.Logger
	a  *auth.JWTAuthenticator
	t  *template.Templates
	ps PostService
	ts TopicService
}

func NewPostHandler(l *slog.Logger, a *auth.JWTAuthenticator,
	t *template.Templates, ps PostService, ts TopicService) *PostHandler {
	return &PostHandler{l: l, a: a, t: t, ps: ps, ts: ts}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := p.ps.GetPostByID(id)
	if err != nil {
		http.Error(rw, "Post Not Found", http.StatusNotFound)
		return
	}

	viewData := new(model.Page)
	viewData.IsAuthor = false

	claims, err := p.a.GetClaimsFromRequest(r)

	if err == nil {
		user := claims["user"].(map[string]interface{})
		userIDFloat := user["id"].(float64)
		userIDInt := int(userIDFloat)

		if post.AuthorId == userIDInt {
			viewData.IsAuthor = true
		}
	}

	data := make(map[string]any)
	data["post"] = post

	viewData.Data = data

	err = p.t.Render(rw, r, "post.page", viewData)
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", err.Error())
		return
	}
}

func (p *PostHandler) GetCreatePost(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		return
	}

	topic, err := p.ts.GetTopicByID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = p.t.Render(rw, r, "create-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", err.Error())
		return
	}
}

func (p *PostHandler) PostCreatePost(rw http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	stringTopicID := r.PostFormValue("topic_id")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		return
	}

	msg := "Failed to get user"

	userValue := r.Context().Value("user")
	if userValue == nil {
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", "cant get value from context")
		return
	}

	user, ok := userValue.(map[string]interface{})
	if !ok {
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", "invalid user type in context")
		return
	}

	userIDFloat, ok := user["id"].(float64)
	if !ok {
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", "invalid user ID type")
		return
	}

	userName, ok := user["name"].(string)
	if !ok {
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", "invalid user name type")
		return
	}

	userID := int(userIDFloat)

	err = p.ps.CreatePost(title, content, id, userID, userName)
	if err != nil {
		msg = "Unable to create post"
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", err.Error())
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%d", id)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := p.ps.GetPostByID(id)
	if err != nil {
		http.Error(rw, "Post Not Found", http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["post"] = post

	err = p.t.Render(rw, r, "edit-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", err.Error())
		return
	}
}

func (p *PostHandler) PostEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	topic, err := p.ts.GetTopicByPostID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	err = p.ps.EditPost(title, content, id)
	if err != nil {
		msg := "Unable to edit post"
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", err.Error())
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetDeletePost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	topic, err := p.ts.GetTopicByPostID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	err = p.ps.DeletePost(id)
	if err != nil {
		msg := "Unable to delete post"
		http.Error(rw, msg, http.StatusInternalServerError)
		p.l.Error(msg, "error", err.Error())
		return
	}

	url := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, url, http.StatusFound)
}
