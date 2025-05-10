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
	query := `SELECT * FROM posts WHERE topic_id = $1`

	rows, err := p.conn.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := new(models.Post)
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.AuthorName,
			&post.TopicId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepository) GetPostByID(postID int) (*models.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`

	post := new(models.Post)

	err := p.conn.QueryRow(query, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorId,
		&post.AuthorName,
		&post.TopicId,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return post, nil
}

func (u *PostRepository) InsertPost(post *models.Post) (int, error) {
	query := `INSERT INTO posts (title, content, topic_id, author_id, author_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := u.conn.QueryRow(query,
		post.Title,
		post.Content,
		post.TopicId,
		post.AuthorId,
		post.AuthorName,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(&post.ID)

	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (u *PostRepository) UpdatePost(post *models.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, topic_id = $3, author_id = $4, author_name = $5, updated_at = $6 WHERE id = $7`

	_, err := u.conn.Exec(query,
		post.Title,
		post.Content,
		post.TopicId,
		post.AuthorId,
		post.AuthorName,
		post.UpdatedAt,
		post.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *PostRepository) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := u.conn.Exec(query, postID)
	if err != nil {
		return err
	}

	return nil
}
