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
