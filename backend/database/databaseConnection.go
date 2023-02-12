package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DBInstance() *mongo.Client {
  err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODO_URI' variable.")
	}

  // Create a new client and connect to the server.
  emptyCtx := context.TODO()	
	client, err := mongo.Connect(emptyCtx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	
	// Ping the primary.
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	
	}
	fmt.Println("Successfully connected and pinged.")	

	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
  coll := client.Database("auth").Collection(collectionName)

	return coll
}