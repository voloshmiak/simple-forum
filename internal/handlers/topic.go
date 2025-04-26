package handlers

import (
	"fmt"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"strconv"
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

func (t *TopicHandler) CreateTopic(rw http.ResponseWriter, r *http.Request) {
	
}

func (t *TopicHandler) UpdateTopic(rw http.ResponseWriter, r *http.Request) {}

func (t *TopicHandler) DeleteTopic(rw http.ResponseWriter, r *http.Request) {}
