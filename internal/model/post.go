package model

import "time"

type Post struct {
	ID         int
	Title      string
	Content    string
	AuthorId   int
	AuthorName string
	TopicId    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
