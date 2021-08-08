package main

import (
	"Sugarbox/models"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func createMappings(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var mappings []models.MovieUserMappingInformation
	_ = json.NewDecoder(request.Body).Decode(&mappings)
	
	for _, v := range mappings {
		result, err := collectionMappings.InsertOne(context.TODO(), v)
		if err != nil {
			fmt.Println("Failed to insert Mapping information; ", err)
			return
		}

		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			return
		}
	}
}
func createUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var users []models.UserInformation
	_ = json.NewDecoder(request.Body).Decode(&users)

	for _, v := range users {
		result, err := collectionUsers.InsertOne(context.TODO(), v)
		if err != nil {
			fmt.Println("Failed to insert Users information for user", v.Name, "; ", err)
			return
		}

		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			return
		}
	}
}
func createMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var movies []models.MovieInformation
	_ = json.NewDecoder(request.Body).Decode(&movies)

	for _, v := range movies {
		result, err := collectionMovies.InsertOne(context.TODO(), v)
		if err != nil {
			fmt.Println("Failed to insert Movies information for movie", v.Name, "; ", err)
			return
		}

		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			return
		}
	}
}
func getUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var users []models.UserInformation

	cur, err := collectionUsers.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Failed to find Users information; ", err)
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.TODO())

	for cur.Next(context.TODO()) {
		var user models.UserInformation
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(writer).Encode(users)
	if err != nil {
		return
	}
}
func getMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var movies []models.MovieInformation

	cur, err := collectionMovies.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Failed to find Movies information; ", err)
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.TODO())

	for cur.Next(context.TODO()) {
		var movie models.MovieInformation
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(writer).Encode(movies)
	if err != nil {
		return
	}
}
func getMappings(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var mappings []models.MovieUserMappingInformation
	cur, err := collectionMappings.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Failed to find Mapping information; ", err)
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.TODO())

	for cur.Next(context.TODO()) {
		var mapping models.MovieUserMappingInformation
		err := cur.Decode(&mapping)
		if err != nil {
			log.Fatal(err)
		}
		mappings = append(mappings, mapping)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(writer).Encode(mappings)
	if err != nil {
		return
	}
}
func deleteCollectionMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	deleteResult, err := collectionMovies.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		return
	}

	err = json.NewEncoder(writer).Encode(deleteResult)
	if err != nil {
		return
	}
}
