package service

import (
	"forum-project/internal/model"
	"time"
)

type PostStorage interface {
	GetPostsByTopicID(topicID int) ([]*model.Post, error)
	GetPostByID(postID int) (*model.Post, error)
	InsertPost(post *model.Post) (int, error)
	UpdatePost(post *model.Post) error
	DeletePost(post *model.Post) error
}

type PostService struct {
	repository PostStorage
}

func NewPostService(repository PostStorage) *PostService {
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

func VerifyPostAuthor(post *model.Post, userID int) bool {
	return post.AuthorId == userID
}

func VerifyPostAuthorOrAdmin(post *model.Post, userID int, userRole string) bool {
	return post.AuthorId == userID || userRole == "admin"
}
