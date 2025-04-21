package models

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorId  string `json:"author_id"`
	TopicId   int    `json:"topic_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewPost() *Post {
	return &Post{}
}
