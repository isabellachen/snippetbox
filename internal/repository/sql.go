package repository

import (
	"database/sql"
	"errors"
	"sync"
	"time"

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

func (m *dbRepo) Create(title string, content string, expires int) (int, error) {
	statement := `INSERT INTO snippets (title, content, created, expires)
								VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? HOUR))`

	result, err := m.db.Exec(statement, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *dbRepo) ById(id int) (*models.Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
								WHERE id = ?`

	row := m.db.QueryRow(statement, id)

	snippet := &models.Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	now := time.Now()

	isExpired := snippet.Expires.Before(now)

	if isExpired {
		return nil, models.ErrIsExpired
	}

	return snippet, nil
}

func (m *dbRepo) Latest(limit int) ([]*models.Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
								WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT ?`

	rows, err := m.db.Query(statement, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		snippet := &models.Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *dbRepo) Last() (*models.Snippet, error) {
	return nil, nil
}
