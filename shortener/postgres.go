package shortener

import (
	"database/sql"
	"errors"
)

// PostgresService is the implementation of service that talks to postgres
type PostgresService struct {
	db *sql.DB
}

// NewPostgresService creates a postgres service
func NewPostgresService(db *sql.DB) (*PostgresService, error) {
	if db == nil {
		return nil, errors.New("empty postgres db")
	}

	s := &PostgresService{
		db: db,
	}

	return s, nil
}

// New creates a link
func (s *PostgresService) New(link *Link) (*Link, error) {
	if link == nil {
		return nil, errors.New("empty link")
	}

	return nil, nil
}

// Get a link info
func (s *PostgresService) Get(slug string) (*Link, error) {
	return nil, nil
}

// Delete a link
func (s *PostgresService) Delete(slug string) error {
	return nil
}
