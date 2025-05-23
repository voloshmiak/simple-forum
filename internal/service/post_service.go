package service

import (
	"forum-project/internal/model"
	"forum-project/internal/repository"
	"time"
)

type PostServicer interface {
	GetPostByID(userID int) (*model.Post, error)
	GetPostsByTopicID(topicID int) ([]*model.Post, error)
	CreatePost(title, content string, topicID, authorID int, authorName string) error
	EditPost(title, content string, postID int) error
	DeletePost(postID int) error
	VerifyPostAuthor(post *model.Post, userID int) bool
	VerifyPostAuthorOrAdmin(post *model.Post, userID int, userRole string) bool
}

type PostService struct {
	repository repository.PostStorage
}

func NewPostService(repository repository.PostStorage) *PostService {
	return &PostService{repository: repository}
}

func (p *PostService) GetPostByID(userID int) (*model.Post, error) {
	post, err := p.repository.GetPostByID(userID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PostService) GetPostsByTopicID(topicID int) ([]*model.Post, error) {
	posts, err := p.repository.GetPostsByTopicID(topicID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostService) CreatePost(title, content string, topicID, authorID int, authorName string) error {
	post := &model.Post{
		Title:      title,
		Content:    content,
		TopicId:    topicID,
		AuthorId:   authorID,
		AuthorName: authorName,
		CreatedAt:  time.Now(),
	}

	postID, err := p.repository.InsertPost(post)
	if err != nil {
		return err
	}

	post.ID = postID
	return nil
}

func (p *PostService) EditPost(title, content string, postID int) error {
	post, err := p.repository.GetPostByID(postID)
	if err != nil {
		return err
	}

	post.Title = title
	post.Content = content

	err = p.repository.UpdatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) DeletePost(postID int) error {
	post, err := p.repository.GetPostByID(postID)
	if err != nil {
		return err
	}
	err = p.repository.DeletePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostService) VerifyPostAuthor(post *model.Post, userID int) bool {
	return post.AuthorId == userID
}

func (p *PostService) VerifyPostAuthorOrAdmin(post *model.Post, userID int, userRole string) bool {
	return post.AuthorId == userID || userRole == "admin"
}
