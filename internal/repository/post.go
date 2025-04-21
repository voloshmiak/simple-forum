package repository

import (
	"database/sql"
	"forum-project/internal/models"
)

type PostRepository struct {
	conn *sql.DB
}

func NewPostRepository(conn *sql.DB) *PostRepository {
	return &PostRepository{conn: conn}
}

func (p *PostRepository) GetPost(postID int) (*models.Post, error) {
	query := `select id, title, content, author_id, topic_id, created_at, updated_at from posts where id = $1`

	post := models.NewPost()

	err := p.conn.QueryRow(query, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorId,
		&post.TopicId,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return post, nil
}
