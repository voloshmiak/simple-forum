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

func (p *PostRepository) GetPostsByTopicID(topicID int) ([]*models.Post, error) {
	query := "SELECT * FROM posts WHERE topic_id = $1"

	rows, err := p.conn.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := models.NewPost()
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.TopicId,
			&post.AuthorId,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepository) GetPostByID(postID int) (*models.Post, error) {
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
