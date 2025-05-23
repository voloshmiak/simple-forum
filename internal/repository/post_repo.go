package repository

import (
	"database/sql"
	"forum-project/internal/model"
)

type PostStorage interface {
	GetPostsByTopicID(topicID int) ([]*model.Post, error)
	GetPostByID(postID int) (*model.Post, error)
	InsertPost(post *model.Post) (int, error)
	UpdatePost(post *model.Post) error
	DeletePost(post *model.Post) error
}

type PostRepository struct {
	conn *sql.DB
}

func NewPostRepository(conn *sql.DB) *PostRepository {
	return &PostRepository{conn: conn}
}

func (p *PostRepository) GetPostsByTopicID(topicID int) ([]*model.Post, error) {
	query := `SELECT * FROM posts WHERE topic_id = $1`

	rows, err := p.conn.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := new(model.Post)
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

func (p *PostRepository) GetPostByID(postID int) (*model.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`

	post := new(model.Post)

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

func (p *PostRepository) InsertPost(post *model.Post) (int, error) {
	query := `INSERT INTO posts (title, content, topic_id, author_id, author_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := p.conn.QueryRow(query,
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

func (p *PostRepository) UpdatePost(post *model.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, topic_id = $3, author_id = $4, author_name = $5, updated_at = $6 WHERE id = $7`

	_, err := p.conn.Exec(query,
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

func (p *PostRepository) DeletePost(post *model.Post) error {
	query := `DELETE FROM posts WHERE id = $1`

	_, err := p.conn.Exec(query, post.ID)
	if err != nil {
		return err
	}

	return nil
}
