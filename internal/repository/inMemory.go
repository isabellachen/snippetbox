package repository

import (
	"fmt"
	"sync"
	"time"

	"snippetbox.isachen.com/internal/models"
)

type inMemoryRepo struct {
	sync.RWMutex
	snippets []models.Snippet
}

func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		snippets: []models.Snippet{},
	}
}

func (r *inMemoryRepo) Create(title string, content string, expires int) (int, error) {
	r.Lock()
	defer r.Unlock()
	s := models.Snippet{}

	s.ID = len(r.snippets) + 1
	s.Title = title
	s.Content = content
	s.Created = time.Now()
	s.Expires = time.Now().Add(time.Hour * time.Duration(expires))
	r.snippets = append(r.snippets, s)
	return s.ID, nil
}

func (r *inMemoryRepo) ById(id int) (*models.Snippet, error) {
	r.RLock()
	defer r.RUnlock()

	s := models.Snippet{}
	if id == 0 {
		return &s, fmt.Errorf("%s: %d", "Invalid ID", id)
	}
	s = r.snippets[id-1]
	return &s, nil
}

func (r *inMemoryRepo) Last() (*models.Snippet, error) {
	r.RLock()
	defer r.RUnlock()

	s := models.Snippet{}
	if len(r.snippets) < 1 {
		return &s, fmt.Errorf("%s", "No snippets available")
	}

	s = r.snippets[len(r.snippets)-1]

	return &s, nil
}
