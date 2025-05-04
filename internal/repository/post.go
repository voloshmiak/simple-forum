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
		post := models.NewPost()
		err := rows.Scan(
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
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *PostRepository) GetPostByID(postID int) (*models.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`

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

func (u *PostRepository) InsertPost(post *models.Post) (int, error) {
	query := `INSERT INTO posts (title, content, topic_id, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := u.conn.QueryRow(query,
		post.Title,
		post.Content,
		post.TopicId,
		post.AuthorId,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(&post.ID)

	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (u *PostRepository) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := u.conn.Exec(query, postID)
	if err != nil {
		return err
	}

	return nil
}
