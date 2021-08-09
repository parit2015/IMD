package main

import (
	"IMD-master/models"
	"IMD-master/utils"
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

func deleteCollections(writer http.ResponseWriter, request *http.Request) {
	var collections = []*mongo.Collection{collectionMovies, collectionUsers, collectionMappings}
	
	for _, v := range collections {
		utils.DeleteDocs(v, writer, request)
	}
}
