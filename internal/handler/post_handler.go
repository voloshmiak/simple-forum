package handler

import (
	"fmt"
	"forum-project/internal/app"
	"forum-project/internal/model"
	"forum-project/internal/service"
	"net/http"
	"strconv"
)

type PostHandler struct {
	app *app.App
}

func NewPostHandler(app *app.App) *PostHandler {
	return &PostHandler{app: app}
}

func (p *PostHandler) GetPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.handleError(rw, "Invalid Post ID", err, http.StatusBadRequest)
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.handleError(rw, "Post Not Found", err, http.StatusNotFound)
		return
	}

	viewData := new(model.Page)
	viewData.IsAuthor = false

	cookie, err := r.Cookie("token")
	if err == nil {
		claims, err := service.ValidateToken(cookie.Value, p.app.Config.JWT.Secret)
		if err == nil {
			user := claims["user"].(map[string]interface{})
			userIDFloat := user["id"].(float64)
			userIDInt := int(userIDFloat)

			isAuthor := service.VerifyPostAuthor(post, userIDInt)
			if isAuthor {
				viewData.IsAuthor = true
			}
		}
	}

	data := make(map[string]any)
	data["post"] = post

	viewData.Data = data

	err = p.app.Templates.Render(rw, r, "post.page", viewData)
	if err != nil {
		p.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (p *PostHandler) GetCreatePost(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		p.handleError(rw, "Invalid Topic ID", err, http.StatusBadRequest)
		return
	}

	topic, err := p.app.TopicService.GetTopicByID(id)
	if err != nil {
		p.handleError(rw, "Topic Not Found", err, http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = p.app.Templates.Render(rw, r, "create-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		p.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (p *PostHandler) PostCreatePost(rw http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	stringTopicID := r.PostFormValue("topic_id")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		p.handleError(rw, "Invalid Topic ID", err, http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*model.AuthorizedUser)
	userID := user.ID
	userName := user.Username

	err = p.app.PostService.CreatePost(title, content, id, userID, userName)
	if err != nil {
		p.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%d", id)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.handleError(rw, "Invalid Post ID", err, http.StatusBadRequest)
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		p.handleError(rw, "Post Not Found", err, http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["post"] = post

	err = p.app.Templates.Render(rw, r, "edit-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		p.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (p *PostHandler) PostEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.handleError(rw, "Invalid Post ID", err, http.StatusBadRequest)
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		p.handleError(rw, "Topic Not Found", err, http.StatusNotFound)
		return
	}

	err = p.app.PostService.EditPost(title, content, id)
	if err != nil {
		p.handleError(rw, "Unable to edit post", err, http.StatusInternalServerError)
		return
	}

	redirectedURL := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, redirectedURL, http.StatusFound)
}

func (p *PostHandler) GetDeletePost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		p.handleError(rw, "Invalid Post ID", err, http.StatusBadRequest)
		return
	}

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		p.handleError(rw, "Topic Not Found", err, http.StatusNotFound)
		return
	}

	err = p.app.PostService.DeletePost(id)
	if err != nil {
		p.handleError(rw, "Unable to delete post", err, http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, url, http.StatusFound)
}

func (p *PostHandler) handleError(rw http.ResponseWriter, msg string, err error, code int) {
	http.Error(rw, msg, code)
	p.app.Logger.Error(msg, "error", err.Error())
}
