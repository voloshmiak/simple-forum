package model

import "time"

type Topic struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	AuthorId    int
}
