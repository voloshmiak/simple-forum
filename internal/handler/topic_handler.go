package handler

import (
	"log/slog"
	"net/http"
	"simple-forum/internal/auth"
	"simple-forum/internal/model"
	"simple-forum/internal/template"
	"strconv"
)

type TopicService interface {
	GetAllTopics() ([]*model.Topic, error)
	GetTopicByID(id int) (*model.Topic, error)
	GetTopicByPostID(id int) (*model.Topic, error)
	CreateTopic(name, description string, authorID int) error
	EditTopic(id int, name, description string) error
	DeleteTopic(id int) error
}

type TopicHandler struct {
	l  *slog.Logger
	a  *auth.JWTAuthenticator
	t  *template.Templates
	ps PostService
	ts TopicService
}

func NewTopicHandler(l *slog.Logger, a *auth.JWTAuthenticator,
	t *template.Templates, ps PostService, ts TopicService) *TopicHandler {
	return &TopicHandler{l: l, a: a, t: t, ps: ps, ts: ts}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.ts.GetAllTopics()
	if err != nil {
		msg := "Unable to get topics"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.t.Render(rw, r, "topics.page", &model.Page{
		Data: data,
	})
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		return
	}

	topic, err := t.ts.GetTopicByID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	posts, err := t.ps.GetPostsByTopicID(id)
	if err != nil {
		msg := "Unable to get posts"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}

	data := make(map[string]any)
	data["posts"] = posts
	data["topic"] = topic

	err = t.t.Render(rw, r, "topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}
}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.t.Render(rw, r, "create-topic.page", new(model.Page))
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}
}

func (t *TopicHandler) PostCreateTopic(rw http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")

	msg := "Failed to get user"

	userValue := r.Context().Value("user")
	if userValue == nil {
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", "cant get value from context")
		return
	}

	user, ok := userValue.(map[string]interface{})
	if !ok {
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", "invalid user type in context")
		return
	}

	userIDFloat, ok := user["id"].(float64)
	if !ok {
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", "invalid user ID type")
		return
	}
	userID := int(userIDFloat)
	err := t.ts.CreateTopic(name, description, userID)
	if err != nil {
		msg = "Unable to create topic"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		return
	}

	topic, err := t.ts.GetTopicByID(id)

	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.t.Render(rw, r, "edit-topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}
}

func (t *TopicHandler) PostEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = t.ts.EditTopic(id, name, description)
	if err != nil {
		msg := "Unable to edit topic"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusBadRequest)
		return
	}

	err = t.ts.DeleteTopic(id)
	if err != nil {
		msg := "Unable to delete topic"
		http.Error(rw, msg, http.StatusInternalServerError)
		t.l.Error(msg, "error", err.Error())
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
