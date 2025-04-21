package service

import (
	"forum-project/internal/models"
	"forum-project/internal/repository"
)

type PostService struct {
	repository *repository.PostRepository
}

func NewPostService(repository *repository.PostRepository) *PostService {
	return &PostService{repository: repository}
}

func (p *PostService) GetPostByID(userID int) (*models.Post, error) {
	post, err := p.repository.GetPostByID(userID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PostService) GetPostsByTopicID(topicID int) ([]*models.Post, error) {
	posts, err := p.repository.GetPostsByTopicID(topicID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
