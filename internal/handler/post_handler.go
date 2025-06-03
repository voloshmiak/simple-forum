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
		http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		http.Error(rw, "Post Not Found", http.StatusNotFound)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
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
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (p *PostHandler) GetCreatePost(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	topic, err := p.app.TopicService.GetTopicByID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = p.app.Templates.Render(rw, r, "create-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
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
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	user := r.Context().Value("user").(*model.AuthorizedUser)
	userID := user.ID
	userName := user.Username

	err = p.app.PostService.CreatePost(title, content, id, userID, userName)
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
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
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	post, err := p.app.PostService.GetPostByID(id)
	if err != nil {
		http.Error(rw, "Post Not Found", http.StatusNotFound)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	data := make(map[string]any)
	data["post"] = post

	err = p.app.Templates.Render(rw, r, "edit-post.page", &model.Page{
		Data: data,
	})
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (p *PostHandler) PostEditPost(rw http.ResponseWriter, r *http.Request) {
	stringPostID := r.PathValue("postID")
	id, err := strconv.Atoi(stringPostID)
	if err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusBadRequest)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	err = p.app.PostService.EditPost(title, content, id)
	if err != nil {
		http.Error(rw, "Unable to edit post", http.StatusInternalServerError)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
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
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	topic, err := p.app.TopicService.GetTopicByPostID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"postID": stringPostID})
		return
	}

	err = p.app.PostService.DeletePost(id)
	if err != nil {
		http.Error(rw, "Unable to delete post", http.StatusInternalServerError)
		p.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}

	url := fmt.Sprintf("/topics/%v", topic.ID)

	http.Redirect(rw, r, url, http.StatusFound)
}
