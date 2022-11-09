package repository

import (
	"database/sql"
	"sync"

	"snippetbox.isachen.com/internal/models"
)

type dbRepo struct {
	sync.RWMutex
	db *sql.DB
}

func NewSqlRepo(dsn string) (*dbRepo, error) {
	db, err := openDB(dsn)
	if err != nil {
		return nil, err
	}
	return &dbRepo{
		db: db,
	}, nil
}

func openDB(dsn string) (*sql.DB, error) {
	// Init a pool of several connections without connecting
	// Uses a imported library driver specific helper
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// This will insert a new snippet into the database.
func (m *dbRepo) Create(title string, content string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.db.Exec(stmt, title, content)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *dbRepo) ById(id int) (*models.Snippet, error) {
	return nil, nil
}

// This will return the 10 most recently created snippets.
func (m *dbRepo) Last() (*models.Snippet, error) {
	return nil, nil
}
