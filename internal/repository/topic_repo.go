package repository

import (
	"database/sql"
	"simple-forum/internal/model"
)

type TopicRepository struct {
	conn *sql.DB
}

func NewTopicRepository(conn *sql.DB) *TopicRepository {
	return &TopicRepository{conn: conn}
}

func (t *TopicRepository) GetAllTopics() ([]*model.Topic, error) {
	query := `SELECT * FROM topics`

	rows, err := t.conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []*model.Topic

	for rows.Next() {
		topic := new(model.Topic)
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

func (t *TopicRepository) GetTopicByID(topicID int) (*model.Topic, error) {
	query := `SELECT * FROM topics WHERE id = $1`

	topic := new(model.Topic)

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

func (t *TopicRepository) GetTopicByPostID(postID int) (*model.Topic, error) {
	query := `SELECT t.* FROM topics t JOIN posts p ON t.id = p.topic_id WHERE p.id = $1`

	topic := new(model.Topic)

	err := t.conn.QueryRow(query, postID).Scan(
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

func (t *TopicRepository) InsertTopic(topic *model.Topic) (int, error) {
	query := `INSERT INTO topics (name, description, created_at, author_id) VALUES ($1, $2, $3, $4) RETURNING id`

	err := t.conn.QueryRow(query,
		topic.Name,
		topic.Description,
		topic.CreatedAt,
		topic.AuthorId).Scan(&topic.ID)

	if err != nil {
		return 0, err
	}

	return topic.ID, nil
}

func (t *TopicRepository) UpdateTopic(topic *model.Topic) error {
	query := `UPDATE topics SET name = $1, description = $2 WHERE id = $3`

	_, err := t.conn.Exec(query,
		topic.Name,
		topic.Description,
		topic.ID)

	if err != nil {
		return err
	}
	return nil
}

func (t *TopicRepository) DeleteTopic(topic *model.Topic) error {
	query := `DELETE FROM topics WHERE id = $1`

	_, err := t.conn.Exec(query, topic.ID)
	if err != nil {
		return err
	}
	return nil
}
