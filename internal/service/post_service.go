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

func (p *PostService) CreatePost(title, content string, topicID, authorID int, authorName string) error {
	post := &models.Post{
		Title:      title,
		Content:    content,
		TopicId:    topicID,
		AuthorId:   authorID,
		AuthorName: authorName,
	}

	postID, err := p.repository.InsertPost(post)
	if err != nil {
		return err
	}

	post.ID = postID
	return nil
}

func (p *PostService) EditPost(post *models.Post) error {
	err := p.repository.UpdatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) DeletePost(postID int) error {
	err := p.repository.DeletePost(postID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) VerifyPostAuthor(post *models.Post, userID int) bool {
	return post.AuthorId == userID
}

func (p *PostService) VerifyPostAuthorOrAdmin(post *models.Post, userID int, userRole string) bool {
	return post.AuthorId == userID || userRole == "admin"
}
