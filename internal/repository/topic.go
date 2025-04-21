package repository

import (
	"database/sql"
	"forum-project/internal/models"
)

type TopicRepository struct {
	conn *sql.DB
}

func NewTopicRepository(conn *sql.DB) *TopicRepository {
	return &TopicRepository{conn: conn}
}

func (t *TopicRepository) GetAllTopics() ([]*models.Topic, error) {
	query := "SELECT id, name, description, created_at, author_id FROM topics"

	rows, err := t.conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []*models.Topic

	for rows.Next() {
		topic := models.NewTopic()
		err = rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Description,
			&topic.CreatedAt,
			&topic.AuthorId,
		)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

func (t *TopicRepository) GetTopicByID(topicID int) (*models.Topic, error) {
	query := `SELECT id, name, description, created_at, author_id FROM topics WHERE id = $1`

	topic := models.NewTopic()

	err := t.conn.QueryRow(query, topicID).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Description,
		&topic.CreatedAt,
		&topic.AuthorId,
	)
	if err != nil {
		return nil, err
	}
	return topic, nil
}
