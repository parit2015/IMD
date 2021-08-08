package main

import (
	"Sugarbox/models"
	"context"
	"encoding/json"
	"fmt"
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
