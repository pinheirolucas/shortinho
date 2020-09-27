package database

import (
	"context"
	"database/sql"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	connectedEngine    Engine
	mongoClient        *mongo.Client
	postgresConnection *sql.DB
)

// CloseFunc is a function to close the database connection
type CloseFunc func() error

// Connect creates a database connection from a given viper config
func Connect(engine Engine, uri string) (CloseFunc, error) {
	if connectedEngine != EngineUnknown {
		return nil, errors.New("database is already initialized")
	}

	switch engine {
	case EngineMongo:
		return connectMongo(uri)
	case EnginePostgres:
		return connectPostgres(uri)
	default:
		return nil, errors.New("unknown database engine")
	}
}

// GetConnection ...
func GetConnection() (Engine, interface{}, error) {
	switch connectedEngine {
	case EngineMongo:
		return connectedEngine, mongoClient, nil
	case EnginePostgres:
		return connectedEngine, postgresConnection, nil
	default:
		return connectedEngine, nil, errors.New("database not initialized")
	}
}

func connectMongo(uri string) (CloseFunc, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(context.Background()); err != nil {
		return nil, err
	}

	connectedEngine = EngineMongo
	mongoClient = client
	return func() error {
		if err := client.Disconnect(context.Background()); err != nil {
			return err
		}

		return nil
	}, nil
}

func connectPostgres(uri string) (CloseFunc, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	connectedEngine = EnginePostgres
	postgresConnection = db
	return db.Close, nil
}
