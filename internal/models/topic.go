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
	CreatedOn   string  `json:"created_on"`
	Posts       []*Post `json:"posts"`
}

var topics = Topics{
	&Topic{
		ID:          1,
		Name:        "Topic 1",
		Description: "Description 1",
		CreatedOn:   "2021-01-01",
		Posts: []*Post{
			{ID: 1, Name: "First post", Content: "Content 1", Author: "Author 1", CreatedOn: "2021-01-01", TopicID: 1},
			{ID: 2, Name: "Question about Topic 1", Content: "I'm having trouble understanding this topic. Can someone explain?", Author: "NewUser", CreatedOn: "2021-01-03", TopicID: 1},
			{ID: 3, Name: "Response to question", Content: "Here's a detailed explanation of Topic 1...", Author: "ExpertUser", CreatedOn: "2021-01-04", TopicID: 1},
		},
	},
	&Topic{
		ID:          2,
		Name:        "Topic 2",
		Description: "Description 2",
		CreatedOn:   "2021-01-02",
		Posts: []*Post{
			{ID: 4, Name: "Starting Topic 2 discussion", Content: "Let's talk about Topic 2 features", Author: "Moderator", CreatedOn: "2021-01-02", TopicID: 2},
			{ID: 5, Name: "My experience", Content: "I've been using this for a month and here's what I found...", Author: "RegularUser", CreatedOn: "2021-01-10", TopicID: 2},
		},
	},
	&Topic{
		ID:          3,
		Name:        "Topic 3",
		Description: "Description 3",
		CreatedOn:   "2021-01-03",
		Posts: []*Post{
			{ID: 6, Name: "Introduction to Topic 3", Content: "Content 2", Author: "Author 2", CreatedOn: "2021-01-02", TopicID: 3},
			{ID: 7, Name: "Topic 3 deep dive", Content: "Here's an in-depth analysis of Topic 3...", Author: "ResearchUser", CreatedOn: "2021-02-05", TopicID: 3},
			{ID: 8, Name: "Recent developments", Content: "Has anyone noticed the recent changes to Topic 3?", Author: "ObservantUser", CreatedOn: "2021-03-10", TopicID: 3},
			{ID: 9, Name: "Help needed", Content: "I'm stuck with a problem related to Topic 3", Author: "NewUser2", CreatedOn: "2021-04-15", TopicID: 3},
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
