package handlers

import (
	"fmt"
	"forum-project/internal/auth"
	"forum-project/internal/models"
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

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.topicService.GetAllTopics()
	if err != nil {
		t.logger.Error(err.Error())
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.templates.Render(rw, r, "topics.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {

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

	data := make(map[string]any)
	data["topic"] = topic

	err = t.templates.Render(rw, r, "topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}

}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.templates.Render(rw, r, "create-topic.page", &models.ViewData{})
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
	stringID := r.PathValue("id")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Topic not found: %s", err))
		http.Error(rw, fmt.Sprintf("Topic not found: %s", err), http.StatusNotFound)
	}

	topic, err := t.topicService.GetTopicByID(topicID)

	if err != nil {
		t.logger.Error(fmt.Sprintf("Topic not found: %s", err))
		http.Error(rw, fmt.Sprintf("Topic not found: %s", err), http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.templates.Render(rw, r, "edit-topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) PostEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Topic not found: %s", err))
		http.Error(rw, fmt.Sprintf("Topic not found: %s", err), http.StatusNotFound)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	_, err = t.topicService.EditTopic(topicID, name, description)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to update topic: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to update topic: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		http.Redirect(rw, r, "/topics/", http.StatusFound)
		return
	}

	topic, err := t.topicService.GetTopicByID(topicID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Topic not found: %s", err))
		http.Error(rw, fmt.Sprintf("Topic not found: %s", err), http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.templates.Render(rw, r, "delete-topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}

func (t *TopicHandler) PostDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		http.Redirect(rw, r, "/topics/", http.StatusFound)
		return
	}

	err = t.topicService.DeleteTopic(id)
	if err != nil {
		t.logger.Error(fmt.Sprintf("Unable to delete topic: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to delete topic: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
