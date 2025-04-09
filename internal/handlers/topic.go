package handlers

import (
	"errors"
	"forum-project/internal/models"
	"forum-project/internal/render"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	logger   *zap.SugaredLogger
	renderer *render.Renderer
}

func NewTopicHandler(logger *zap.SugaredLogger, renderer *render.Renderer) *TopicHandler {
	return &TopicHandler{logger, renderer}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	t.logger.Info("Handle GET Topics")

	topics := models.GetTopics()

	t.renderer.RenderTemplate(rw, "topics.page", topics)
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	t.logger.Info("Handle GET Topic")
	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)

	if err != nil {
		http.Redirect(rw, r, "/topics/", http.StatusFound)
		return
	}

	topic, err := models.FindTopic(id)
	if errors.Is(err, models.TopicNotFoundError) {
		t.logger.Error("Topic not found", err)
		http.Error(rw, "Topic not found", http.StatusNotFound)
		return
	}

	t.renderer.RenderTemplate(rw, "topic.page", topic)

}
