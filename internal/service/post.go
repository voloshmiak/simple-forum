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

func (p *PostService) CreatePost(title, content string, topicID, authorID int) (*models.Post, error) {
	post := models.NewPost()
	post.Title = title
	post.Content = content
	post.CreatedAt = "Now"
	post.UpdatedAt = "Now"
	post.TopicId = topicID
	post.AuthorId = authorID

	postID, err := p.repository.InsertPost(post)
	if err != nil {
		return nil, err
	}

	post.ID = postID
	return post, nil
}
