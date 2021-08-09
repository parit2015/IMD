package main

import (
	"IMD-master/models"
	"IMD-master/utils"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func createMappings(writer http.ResponseWriter, request *http.Request) {
	var mappings []models.MovieUserMappingInformation
	utils.InsertDocs(collectionMappings, writer, request, mappings)
}

func createUsers(writer http.ResponseWriter, request *http.Request) {
	var users []models.UserInformation
	utils.InsertDocs(collectionUsers, writer, request, users)
}

func createMovie(writer http.ResponseWriter, request *http.Request) {
	var movies []models.MovieInformation
	utils.InsertDocs(collectionMovies, writer, request, movies)
}

func deleteDocs(collection *mongo.Collection, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	deleteResult, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		return
	}

	err = json.NewEncoder(writer).Encode(deleteResult)
	if err != nil {
		return
	}
}

func deleteCollections(writer http.ResponseWriter, request *http.Request) {
	var collections = []*mongo.Collection{collectionMovies, collectionUsers, collectionMappings}
	
	for _, v := range collections {
		deleteDocs(v, writer, request)
	}
}
