package shortener

import (
	"database/sql"
	"errors"

	"github.com/pinheirolucas/shortinho/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// Link relates a slug to a target url
type Link struct{}

// Service for managing links
type Service interface {
	New(link *Link) (*Link, error)
	Get(slug string) (*Link, error)
	Delete(slug string) error
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
