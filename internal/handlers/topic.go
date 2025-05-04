package handlers

import (
	"fmt"
	"forum-project/internal/auth"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type TopicHandler struct {
	logger       *slog.Logger
	templates    *template.Manager
	topicService *service.TopicService
}

func NewTopicHandler(logger *slog.Logger, renderer *template.Manager, topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{logger, renderer, topicService}
}

func (t *TopicHandler) GetAllTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.topicService.GetAllTopics()
	if err != nil {
		t.logger.Error(err.Error())
	}

	err = t.templates.Render(rw, "topics.page", topics)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) GetTopicByID(rw http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)

	if err != nil {
		http.Redirect(rw, r, "/topics/", http.StatusFound)
		return
	}

	topic, err := t.topicService.GetTopicByID(id)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Topic not found: %s", err))
		http.Error(rw, fmt.Sprintf("Topic not found: %s", err), http.StatusNotFound)
		return
	}

	err = t.templates.Render(rw, "topic.page", topic)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}

}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.templates.Render(rw, "create-topic.page", nil)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) PostCreateTopic(rw http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	authorIDfloat, ok := claims["id"].(float64)
	if !ok {
		t.logger.Error("Unable to convert author id to float64")
		http.Error(rw, "Unable to convert author id to float64", http.StatusBadRequest)
		return
	}
	authorIDInt := int(authorIDfloat)

	_, err = t.topicService.CreateTopic(name, description, authorIDInt)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to create topic: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to create topic: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetEditTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.templates.Render(rw, "create-topic.page", nil)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) PutEditTopic(rw http.ResponseWriter, r *http.Request) {}

func (t *TopicHandler) GetRemoveTopic(rw http.ResponseWriter, r *http.Request) {}

func (t *TopicHandler) DeleteRemoveTopic(rw http.ResponseWriter, r *http.Request) {}
