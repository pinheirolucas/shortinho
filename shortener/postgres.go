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

	_, err := s.Get(link.Slug)
	switch err {
	case ErrNotFound:
		// continue
	case nil:
		return nil, ErrAlreadyExists
	default:
		return nil, err
	}

	var maxHits sql.NullInt64
	if link.MaxHits != 0 {
		maxHits.Valid = true
		maxHits.Int64 = int64(link.MaxHits)
	}

	_, err = s.db.Exec(`
		insert into link (
			slug,
			target_url,
			max_hits
		) values (
			$1,
			$2,
			$3
		)
	`, link.Slug, link.TargetURL, maxHits)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// Get a link info
func (s *PostgresService) Get(slug string) (*Link, error) {
	var (
		dbSlug    string
		targetURL string
		active    bool
		hits      sql.NullInt64
		maxHits   sql.NullInt64
	)

	err := s.db.QueryRow(`
		select
			l.slug,
			l.target_url,
			l.max_hits,
			l.active,
			count(h)
		from
			link l
		left join link_hit h on
			h.slug = l.slug and
			not h.deleted
		where
			l.slug = $1
		group by
			l.slug,
			l.target_url,
			l.max_hits,
			l.active;
	`, slug).Scan(
		&dbSlug,
		&targetURL,
		&maxHits,
		&active,
		&hits,
	)
	switch err {
	case nil:
		// continue
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}

	link := &Link{
		Slug:      dbSlug,
		TargetURL: targetURL,
		Active:    active,
	}

	if hits.Valid {
		link.Hits = int(hits.Int64)
	}

	if maxHits.Valid {
		link.MaxHits = int(maxHits.Int64)
	}

	return link, nil
}

// Activate a link
func (s *PostgresService) Activate(slug string) error {
	return s.setActive(slug, true)
}

// Deactivate a link
func (s *PostgresService) Deactivate(slug string) error {
	return s.setActive(slug, false)
}

func (s *PostgresService) setActive(slug string, active bool) error {
	result, err := s.db.Exec(`
		update link
		set active = $1
		where slug = $2
	`, active, slug)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return ErrNotFound
	}

	return nil
}

// Hit the link
func (s *PostgresService) Hit(slug string) error {
	_, err := s.db.Exec("insert into link_hit (slug) values ($1)", slug)
	if err != nil {
		return err
	}

	return nil
}

// Reset link stats
func (s *PostgresService) Reset(slug string) error {
	result, err := s.db.Exec(`
		update link_hit
		set deleted = true
		where slug = $1
	`, slug)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return ErrNotFound
	}

	return nil
}
