package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


const uri = "mongodb://localhost:27017"

func ConnectDatabase() *mongo.Client{

	log.Println("Connecting to MongoDB....")
	
	ctx , cancel := context.WithTimeout(context.Background() , 10 * time.Second)
	defer cancel()

	client , err := mongo.Connect(ctx , options.Client().ApplyURI(uri))

	if err != nil{
		log.Fatalf("Failed to connect to MongoDB: %v" , err)
	}

	err = client.Ping(ctx , nil)

	if err != nil {
		log.Fatalf("MongoDB ping failed: %v" , err)
	}

	log.Println("Successfully connected to MongoDB!")
	
	return client
}

var client *mongo.Client = ConnectDatabase()

func OpenCollection(collectionName string) *mongo.Collection {
	if client == nil{
		log.Fatal("MongoDB Client is not initialized. Please connect DB first")
	}

	return client.Database("1to1Chat").Collection(collectionName)
}