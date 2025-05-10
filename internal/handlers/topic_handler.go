package handlers

import (
	"forum-project/internal/config"
	"forum-project/internal/models"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	app *config.AppConfig
}

func NewTopicHandler(app *config.AppConfig) *TopicHandler {
	return &TopicHandler{app: app}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.app.TopicService.GetAllTopics()
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to get topics", err)
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.app.Templates.Render(rw, r, "topics.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.app.Errors.BadRequest(rw, "Invalid Topic ID", err)
		return
	}

	topic, err := t.app.TopicService.GetTopicByID(id)
	if err != nil {
		t.app.Errors.NotFound(rw, "Topic Not Found", err)
		return
	}

	posts, err := t.app.PostService.GetPostsByTopicID(id)
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to get posts", err)
		return
	}

	data := make(map[string]any)
	data["posts"] = posts
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.app.Templates.Render(rw, r, "create-topic.page", &models.ViewData{})
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) PostCreateTopic(rw http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")

	user := r.Context().Value("user").(*models.AuthorizedUser)
	userID := user.ID

	err := t.app.TopicService.CreateTopic(name, description, userID)
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to create topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.app.Errors.BadRequest(rw, "Invalid Topic ID", err)
	}

	topic, err := t.app.TopicService.GetTopicByID(id)

	if err != nil {
		t.app.Errors.NotFound(rw, "Topic Not Found", err)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "edit-topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) PostEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.app.Errors.BadRequest(rw, "Invalid Topic ID", err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = t.app.TopicService.EditTopic(id, name, description)
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to edit topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.app.Errors.BadRequest(rw, "Invalid Topic ID", err)
		return
	}

	err = t.app.TopicService.DeleteTopic(id)
	if err != nil {
		t.app.Errors.InternalServer(rw, "Unable to delete topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
