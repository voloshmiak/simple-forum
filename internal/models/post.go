package models

import "time"

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorId   int       `json:"author_id"`
	AuthorName string    `json:"author_name"`
	TopicId    int       `json:"topic_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewPost() *Post {
	return &Post{}
}
