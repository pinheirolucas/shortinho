package shortener

import (
	"database/sql"
	"errors"

	"github.com/pinheirolucas/shortinho/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service errors
var (
	ErrAlreadyExists = errors.New("link already exists")
	ErrNotFound      = errors.New("link not found")
)

// Link relates a slug to a target url
type Link struct {
	Slug      string `json:"slug,omitempty"`
	TargetURL string `json:"targetUrl,omitempty"`
	Hits      int    `json:"hits,omitempty"`
	MaxHits   int    `json:"maxHits,omitempty"`
	Active    bool   `json:"active,omitempty"`
}

// Service for managing links
type Service interface {
	New(link *Link) (*Link, error)
	Get(slug string) (*Link, error)
	Activate(slug string) error
	Deactivate(slug string) error
	Hit(slug string) error
	Reset(slug string) error
}

// NewService creates a new link service from the loaded database engine
func NewService() (Service, error) {
	engine, connection, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	switch engine {
	case database.EngineMongo:
		return NewMongoService(connection.(*mongo.Client))
	case database.EnginePostgres:
		return NewPostgresService(connection.(*sql.DB))
	}

	return nil, errors.New("invalid initialized database")
}
