package models

import (
	"time"
)

type Snippet struct {
	ID      int
	Title   string `json:"title"`
	Content string `json:"content"`
	Created time.Time
	Expires time.Time
}

type Repository interface {
	Create(title string, content string) (int, error)
	ById(id int) (*Snippet, error)
	Last() (*Snippet, error)
}
