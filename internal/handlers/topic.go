package handlers

import (
	"forum-project/internal/models"
	"forum-project/internal/mylogger"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	logger       *mylogger.WrappedLogger
	templates    *template.Manager
	topicService *service.TopicService
	postService  *service.PostService
}

func NewTopicHandler(logger *mylogger.WrappedLogger, renderer *template.Manager, topicService *service.TopicService, postService *service.PostService) *TopicHandler {
	return &TopicHandler{logger, renderer, topicService, postService}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.topicService.GetAllTopics()
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to get topics", err)
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.templates.Render(rw, r, "topics.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) GetTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.logger.BadRequestError(rw, "Invalid Topic ID", err)
		return
	}

	topic, err := t.topicService.GetTopicByID(id)
	if err != nil {
		t.logger.NotFoundError(rw, "Topic Not Found", err)
		return
	}

	posts, err := t.postService.GetPostsByTopicID(id)
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to get posts", err)
		return
	}

	data := make(map[string]any)
	data["posts"] = posts
	data["topic"] = topic

	err = t.templates.Render(rw, r, "topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.templates.Render(rw, r, "create-topic.page", &models.ViewData{})
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) PostCreateTopic(rw http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")

	user := r.Context().Value("user")
	userIDfloat := user.(map[string]interface{})["id"].(float64)
	userID := int(userIDfloat)

	err := t.topicService.CreateTopic(name, description, userID)
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to create topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.logger.BadRequestError(rw, "Invalid Topic ID", err)
	}

	topic, err := t.topicService.GetTopicByID(id)

	if err != nil {
		t.logger.NotFoundError(rw, "Topic Not Found", err)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.templates.Render(rw, r, "edit-topic.page", &models.ViewData{
		Data: data,
	})
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}

func (t *TopicHandler) PostEditTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.logger.BadRequestError(rw, "Invalid Topic ID", err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	err = t.topicService.EditTopic(id, name, description)
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to edit topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}

func (t *TopicHandler) GetDeleteTopic(rw http.ResponseWriter, r *http.Request) {
	stringTopicID := r.PathValue("topicID")
	id, err := strconv.Atoi(stringTopicID)
	if err != nil {
		t.logger.BadRequestError(rw, "Invalid Topic ID", err)
		return
	}

	err = t.topicService.DeleteTopic(id)
	if err != nil {
		t.logger.ServerInternalError(rw, "Unable to delete topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
