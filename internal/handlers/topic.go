package handlers

import (
	"errors"
	"fmt"
	"forum-project/internal/models"
	"forum-project/internal/render"
	"log/slog"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	logger   *slog.Logger
	renderer *render.Renderer
}

func NewTopicHandler(logger *slog.Logger, renderer *render.Renderer) *TopicHandler {
	return &TopicHandler{logger, renderer}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {

	topics := models.GetTopics()

	err := t.renderer.RenderTemplate(rw, "topics.page", topics)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to render template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to render template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {

	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)

	if err != nil {
		http.Redirect(rw, r, "/topics/", http.StatusFound)
		return
	}

	topic, err := models.FindTopic(id)
	if errors.Is(err, models.TopicNotFoundError) {
		t.logger.Error(fmt.Sprintf("Topic not found: %s", err))
		http.Error(rw, fmt.Sprintf("Topic not found: %s", err), http.StatusNotFound)
		return
	}

	err = t.renderer.RenderTemplate(rw, "topic.page", topic)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to render template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to render template: %s", err), http.StatusInternalServerError)
	}

}

func (t *TopicHandler) CreateTopic(rw http.ResponseWriter, r *http.Request) {}

func (t *TopicHandler) UpdateTopic(rw http.ResponseWriter, r *http.Request) {}

func (t *TopicHandler) DeleteTopic(rw http.ResponseWriter, r *http.Request) {}
