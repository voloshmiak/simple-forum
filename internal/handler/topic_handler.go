package handler

import (
	"forum-project/internal/app"
	"forum-project/internal/model"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	app *app.App
}

func NewTopicHandler(app *app.App) *TopicHandler {
	return &TopicHandler{app: app}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.app.TopicService.GetAllTopics()
	if err != nil {
		http.Error(rw, "Unable to get topics", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.app.Templates.Render(rw, r, "topics.page", &model.Page{
		Data: data,
	})
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	topic, err := t.app.TopicService.GetTopicByID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	posts, err := t.app.PostService.GetPostsByTopicID(id)
	if err != nil {
		http.Error(rw, "Unable to get posts", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "path", r.URL.Path)
		return
	}

	data := make(map[string]any)
	data["posts"] = posts
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.app.Templates.Render(rw, r, "create-topic.page", new(model.Page))
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (t *TopicHandler) PostCreateTopic(rw http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")

	user := r.Context().Value("user").(*model.AuthorizedUser)
	userID := user.ID

	err := t.app.TopicService.CreateTopic(name, description, userID)
	if err != nil {
		http.Error(rw, "Unable to create topic", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	topic, err := t.app.TopicService.GetTopicByID(id)

	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusNotFound, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "edit-topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (t *TopicHandler) PostEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = t.app.TopicService.EditTopic(id, name, description)
	if err != nil {
		http.Error(rw, "Unable to edit topic", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "path", r.URL.Path)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		t.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusBadRequest, "path", r.URL.Path, "context", map[string]interface{}{"topicID": stringTopicID})
		return
	}

	err = t.app.TopicService.DeleteTopic(id)
	if err != nil {
		http.Error(rw, "Unable to delete topic", http.StatusInternalServerError)
		t.app.Logger.Error(err.Error(), "method", r.Method, "path", r.URL.Path)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
