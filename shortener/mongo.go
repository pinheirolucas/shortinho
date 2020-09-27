package shortener

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoService is the implementation of service that talks to postgres
type MongoService struct {
	client *mongo.Client
}

// NewMongoService creates a postgres service
func NewMongoService(client *mongo.Client) (*MongoService, error) {
	if client == nil {
		return nil, errors.New("empty mongo client")
	}

	s := &MongoService{
		client: client,
	}

	return s, nil
}

// New creates a link
func (s *MongoService) New(link *Link) (*Link, error) {
	return nil, errors.New("not implemented")
}

// Get a link info
func (s *MongoService) Get(slug string) (*Link, error) {
	return nil, errors.New("not implemented")
}

// Activate a link
func (s *MongoService) Activate(slug string) error {
	return nil
}

// Deactivate a link
func (s *MongoService) Deactivate(slug string) error {
	return nil
}

// Hit the link
func (s *MongoService) Hit(slug string) error {
	return nil
}

// Reset link stats
func (s *MongoService) Reset(slug string) error {
	return nil
}
