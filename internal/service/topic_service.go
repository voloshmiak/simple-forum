package service

import (
	"forum-project/internal/models"
	"forum-project/internal/repository"
)

type TopicService struct {
	repository *repository.TopicRepository
}

func NewTopicService(repository *repository.TopicRepository) *TopicService {
	return &TopicService{repository: repository}
}

func (t *TopicService) GetAllTopics() ([]*models.Topic, error) {
	topics, err := t.repository.GetAllTopics()
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func (t *TopicService) GetTopicByID(id int) (*models.Topic, error) {
	topic, err := t.repository.GetTopicByID(id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (t *TopicService) GetTopicByPostID(id int) (*models.Topic, error) {
	topic, err := t.repository.GetTopicByPostID(id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (t *TopicService) CreateTopic(name, description string, authorID int) error {
	topic := &models.Topic{
		Name:        name,
		Description: description,
		AuthorId:    authorID,
	}

	_, err := t.repository.InsertTopic(topic)
	if err != nil {
		return err
	}

	return nil
}

func (t *TopicService) EditTopic(id int, name, description string) error {
	topic := &models.Topic{
		ID:          id,
		Name:        name,
		Description: description,
	}

	err := t.repository.UpdateTopic(topic)
	if err != nil {
		return err
	}
	return nil
}

func (t *TopicService) DeleteTopic(id int) error {
	err := t.repository.DeleteTopic(id)
	if err != nil {
		return err
	}
	return nil
}
