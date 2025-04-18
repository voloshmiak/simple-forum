package models

import (
	"encoding/json"
	"fmt"
	"io"
)

type Topic struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   string  `json:"created_by"`
	Posts       []*Post `json:"posts"`
}

var topics = Topics{
	&Topic{
		ID:          1,
		Name:        "Topic 1",
		Description: "Description 1",
		CreatedAt:   "2021-01-01",
		Posts: []*Post{
			{ID: 1, Title: "First post", Content: "Content 1", AuthorId: "Author 1", CreatedAt: "2021-01-01", TopicId: 1},
			{ID: 2, Title: "Question about Topic 1", Content: "I'm having trouble understanding this topic. Can someone explain?", AuthorId: "NewUser", CreatedAt: "2021-01-03", TopicId: 1},
			{ID: 3, Title: "Response to question", Content: "Here's a detailed explanation of Topic 1...", AuthorId: "ExpertUser", CreatedAt: "2021-01-04", TopicId: 1},
		},
	},
	&Topic{
		ID:          2,
		Name:        "Topic 2",
		Description: "Description 2",
		CreatedAt:   "2021-01-02",
		Posts: []*Post{
			{ID: 4, Title: "Starting Topic 2 discussion", Content: "Let's talk about Topic 2 features", AuthorId: "Moderator", CreatedAt: "2021-01-02", TopicId: 2},
			{ID: 5, Title: "My experience", Content: "I've been using this for a month and here's what I found...", AuthorId: "RegularUser", CreatedAt: "2021-01-10", TopicId: 2},
		},
	},
	&Topic{
		ID:          3,
		Name:        "Topic 3",
		Description: "Description 3",
		CreatedAt:   "2021-01-03",
		Posts: []*Post{
			{ID: 6, Title: "Introduction to Topic 3", Content: "Content 2", AuthorId: "Author 2", CreatedAt: "2021-01-02", TopicId: 3},
			{ID: 7, Title: "Topic 3 deep dive", Content: "Here's an in-depth analysis of Topic 3...", AuthorId: "ResearchUser", CreatedAt: "2021-02-05", TopicId: 3},
			{ID: 8, Title: "Recent developments", Content: "Has anyone noticed the recent changes to Topic 3?", AuthorId: "ObservantUser", CreatedAt: "2021-03-10", TopicId: 3},
			{ID: 9, Title: "Help needed", Content: "I'm stuck with a problem related to Topic 3", AuthorId: "NewUser2", CreatedAt: "2021-04-15", TopicId: 3},
		},
	},
}

type Topics []*Topic

func (t *Topics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Topic) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Topics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(t)
}

func GetTopics() Topics {
	return topics
}

var TopicNotFoundError = fmt.Errorf("topic not found error")

func FindTopic(id int) (*Topic, error) {
	for _, t := range topics {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, TopicNotFoundError
}
