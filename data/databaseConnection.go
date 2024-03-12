package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Client       *mongo.Client
	DatabaseName string
}

func NewMongoDBinstance() *MongoConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("[ERROR] loading .env file")
	}

	MongoDB := os.Getenv("DB_MONGO_URL")
	DatabaseName := os.Getenv("MONGO_DATABASE_NAME")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return &MongoConfig{
		Client:       client,
		DatabaseName: DatabaseName,
	}
}

var MongoCfg = NewMongoDBinstance()
var Client *mongo.Client = MongoCfg.Client

func OpenCollection(client *mongo.Client, DatabaseName, collectionName string) *mongo.Collection {
	fmt.Println("Connected to", DatabaseName)
	var collection *mongo.Collection = client.Database(DatabaseName).Collection(collectionName)
	return collection
}
