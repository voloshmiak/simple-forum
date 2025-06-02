package handler

import (
	"forum-project/internal/application"
	"forum-project/internal/model"
	"net/http"
	"strconv"
)

type TopicHandler struct {
	app *application.App
}

func NewTopicHandler(app *application.App) *TopicHandler {
	return &TopicHandler{app: app}
}

func (t *TopicHandler) GetTopics(rw http.ResponseWriter, r *http.Request) {
	topics, err := t.app.TopicService.GetAllTopics()
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to get topics", err)
		return
	}

	data := make(map[string]any)
	data["topics"] = topics

	err = t.app.Templates.Render(rw, r, "topics.page", &model.Page{
		Data: data,
	})
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to render template", err)
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

	topic, err := t.app.TopicService.GetTopicByID(id)
	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	posts, err := t.app.PostService.GetPostsByTopicID(id)
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to get posts", err)
		return
	}

	var postRows [][]model.Post
	for i := 0; i < len(posts); i += 2 {
		row := []model.Post{*posts[i]}
		if i+1 < len(posts) {
			row = append(row, *posts[i+1])
		}
		postRows = append(postRows, row)
	}

	data := make(map[string]any)
	data["posts"] = postRows
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to render template", err)
		return
	}
}

func (t *TopicHandler) GetCreateTopic(rw http.ResponseWriter, r *http.Request) {
	err := t.app.Templates.Render(rw, r, "create-topic.page", new(model.Page))
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to render template", err)
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
		InternalServer(t.app.Logger, rw, "Unable to create topic", err)
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

	topic, err := t.app.TopicService.GetTopicByID(id)

	if err != nil {
		http.Error(rw, "Topic Not Found", http.StatusNotFound)
		return
	}

	data := make(map[string]any)
	data["topic"] = topic

	err = t.app.Templates.Render(rw, r, "edit-topic.page", &model.Page{
		Data: data,
	})
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to render template", err)
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

	err = t.app.TopicService.EditTopic(id, name, description)
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to edit topic", err)
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

	err = t.app.TopicService.DeleteTopic(id)
	if err != nil {
		InternalServer(t.app.Logger, rw, "Unable to delete topic", err)
		return
	}

	http.Redirect(rw, r, "/topics", http.StatusFound)
}
