package database

import (
	"context"
	"os"

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
		panic(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	DB = client
	return client, nil
}

func Start() error {
	_, err := Connect()
	if err != nil {
		return err
	}
	return nil
}

func GetDB() *mongo.Client {
	return DB
}
