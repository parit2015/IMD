package helper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func GetDBConnection(uri string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Unable to get db client; ", err)
		os.Exit(1)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}
 
func GetCollection(client *mongo.Client) (*mongo.Collection, *mongo.Collection, *mongo.Collection) {
	collectionMovies := client.Database("IMDB").Collection("movies")
	collectionUsers := client.Database("IMDB").Collection("users")
	collectionMappings := client.Database("IMDB").Collection("mappings")

	return collectionMovies, collectionUsers, collectionMappings
}