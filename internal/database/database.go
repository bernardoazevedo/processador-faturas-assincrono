package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var DB *mongo.Client

func Connect() (*mongo.Client, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://" + user + ":" + pass + "@mongodb:27017"))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	DB = client
	return client, nil
}

func Start() error {
	var lastErr error

	for range 5 {
		_, err := Connect()
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(2 * time.Second)
	}

	return errors.New("failed to connect to MongoDB after 5 attempts: " + lastErr.Error())
}

func GetDB() *mongo.Client {
	return DB
}
