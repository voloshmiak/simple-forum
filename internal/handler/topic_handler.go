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
		t.handleError(rw, "Unable to get topics", err, http.StatusInternalServerError)
		return
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.app.Templates.Render(rw, r, "topics.page", &model.Page{
		Data: data,
	})
	if err != nil {
		t.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.handleError(rw, "Invalid Topic ID", err, http.StatusBadRequest)
		return
	}

	topic, err := t.app.TopicService.GetTopicByID(id)
	if err != nil {
		t.handleError(rw, "Topic Not Found", err, http.StatusNotFound)
		return
	}

	posts, err := t.app.PostService.GetPostsByTopicID(id)
	if err != nil {
		t.handleError(rw, "Unable to get posts", err, http.StatusInternalServerError)
		return
	}

	data := make(map[string]any)
	data["posts"] = posts
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		t.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.app.Templates.Render(rw, r, "create-topic.page", new(model.Page))
	if err != nil {
		t.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
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
		t.handleError(rw, "Unable to create topic", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.handleError(rw, "Invalid Topic ID", err, http.StatusBadRequest)
		return
	}

	topic, err := t.app.TopicService.GetTopicByID(id)

	if err != nil {
		t.handleError(rw, "Topic Not Found", err, http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "edit-topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		t.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (t *TopicHandler) PostEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.handleError(rw, "Invalid Topic ID", err, http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = t.app.TopicService.EditTopic(id, name, description)
	if err != nil {
		t.handleError(rw, "Unable to edit topic", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.handleError(rw, "Invalid Topic ID", err, http.StatusBadRequest)
		return
	}

	err = t.app.TopicService.DeleteTopic(id)
	if err != nil {
		t.handleError(rw, "Unable to delete topic", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) handleError(rw http.ResponseWriter, msg string, err error, code int) {
	http.Error(rw, msg, code)
	t.app.Logger.Error(msg, "error", err.Error())
}
