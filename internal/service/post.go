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

func (p *PostService) GetPost(userID int) (*models.Post, error) {
	post, err := p.repository.GetPost(userID)
	if err != nil {
		return nil, err
	}
	return post, nil
}
