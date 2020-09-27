package shortener

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	linkCollection := s.client.Database("shortinho").Collection("link")
	link.Active = true
	link.CreatedAt = int(time.Now().UnixNano() / 1e6)
	link.UpdatedAt = int(time.Now().UnixNano() / 1e6)

	_, err = linkCollection.InsertOne(context.Background(), link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// Get a link info
func (s *MongoService) Get(slug string) (*Link, error) {
	link := new(Link)
	linkCollection := s.client.Database("shortinho").Collection("link")

	err := linkCollection.FindOne(context.Background(), bson.M{
		"slug": slug,
	}).Decode(link)
	switch err {
	case nil:
		// continue
	case mongo.ErrNoDocuments:
		return nil, ErrNotFound
	default:
		return nil, err
	}

	return link, nil
}

// Activate a link
func (s *MongoService) Activate(slug string) error {
	return s.setActive(slug, true)
}

// Deactivate a link
func (s *MongoService) Deactivate(slug string) error {
	return s.setActive(slug, false)
}

func (s *MongoService) setActive(slug string, active bool) error {
	linkCollection := s.client.Database("shortinho").Collection("link")
	_, err := linkCollection.UpdateOne(
		context.Background(),
		bson.M{"slug": slug},
		bson.M{
			"$set": bson.M{
				"active": active,
			},
		},
	)
	if err != nil {
		return err
	}

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
