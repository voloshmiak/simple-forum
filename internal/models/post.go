package models

import (
	"encoding/json"
	"fmt"
	"io"
)

type Post struct {
	ID        int    `json:"id"`
	Name      string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	TopicID   int    `json:"topic_id"`
	CreatedOn string `json:"created_on"`
}

type Posts []*Post

func (p *Posts) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Posts) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Post) ToJSON(w io.Writer) error {
	d := json.NewEncoder(w)
	return d.Encode(p)
}

func GetTopicPosts(topicID int) (Posts, error) {
	topic, err := FindTopic(topicID)
	if err != nil {
		return nil, err
	}

	return topic.Posts, nil

}

var PostNotFoundError = fmt.Errorf("post not found error")

func FindPost(postID int) (*Post, error) {
	for _, topic := range topics {
		for _, post := range topic.Posts {
			if post.ID == postID {
				return post, nil
			}
		}
	}

	return nil, PostNotFoundError
}
