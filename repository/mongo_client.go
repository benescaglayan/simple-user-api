package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

const (
	DatabaseName = "Company"
)

func InitDatabase(databaseUri string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(databaseUri))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Println(err)
		panic(err)
	}

	database := client.Database(DatabaseName)

	return database
}
