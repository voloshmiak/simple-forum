package service

import (
	"simple-forum/internal/model"
	"time"
)

type TopicStorage interface {
	GetAllTopics() ([]*model.Topic, error)
	GetTopicByID(topicID int) (*model.Topic, error)
	GetTopicByPostID(postID int) (*model.Topic, error)
	InsertTopic(topic *model.Topic) (int, error)
	UpdateTopic(topic *model.Topic) error
	DeleteTopic(topic *model.Topic) error
}

type TopicService struct {
	repository TopicStorage
}

func NewTopicService(repository TopicStorage) *TopicService {
	return &TopicService{repository: repository}
}

func (t *TopicService) GetAllTopics() ([]*model.Topic, error) {
	topics, err := t.repository.GetAllTopics()
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func (t *TopicService) GetTopicByID(id int) (*model.Topic, error) {
	topic, err := t.repository.GetTopicByID(id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (t *TopicService) GetTopicByPostID(id int) (*model.Topic, error) {
	topic, err := t.repository.GetTopicByPostID(id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (t *TopicService) CreateTopic(name, description string, authorID int) error {
	topic := &model.Topic{
		Name:        name,
		Description: description,
		AuthorId:    authorID,
		CreatedAt:   time.Now(),
	}

	_, err := t.repository.InsertTopic(topic)
	if err != nil {
		return err
	}

	return nil
}

func (t *TopicService) EditTopic(id int, name, description string) error {
	topic := &model.Topic{
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
	topic, err := t.repository.GetTopicByID(id)
	if err != nil {
		return err
	}
	err = t.repository.DeleteTopic(topic)
	if err != nil {
		return err
	}
	return nil
}
