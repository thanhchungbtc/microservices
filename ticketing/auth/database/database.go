package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
}

func New() (*Database, error) {

	db, err := connectToDB()
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func connectToDB() (*mongo.Database, error) {
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://auth-mongo-srv:27017"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}
	db := client.Database("users")
	return db, nil
}
