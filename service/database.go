package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB *mongo.Database
)

func ConnectMongoDB() {
	uri := os.Getenv("MONGO_URI")
	credentials := options.Credential{
		Username: os.Getenv("MONGO_ROOT_USERNAME"),
		Password: os.Getenv("MONGO_ROOT_PASSWORD"),
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("ReportSystem")
	fmt.Println("Successfully connected to MongoDB")
}
