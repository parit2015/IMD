package main

import (
	"IMD-master/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func getOldNewDocument(request *http.Request, collection *mongo.Collection,
	DocOld *models.MovieUserMappingInformation, DocNew *models.MovieUserMappingInformation, filter bson.M) {

	findOne(collection, filter, DocOld)

	_ = json.NewDecoder(request.Body).Decode(DocNew)
}

func updateDB(collection *mongo.Collection, filter bson.M, tobeUpdateBSONDoc bson.D,
	targetObj *models.MovieUserMappingInformation) {

	findOneAndUpdate(collection, filter, tobeUpdateBSONDoc, targetObj)
}

func createMappings(writer http.ResponseWriter, request *http.Request) {
	var mappings []models.MovieUserMappingInformation
	insertDocs(collectionMappings, writer, request, mappings)
}

func createUsers(writer http.ResponseWriter, request *http.Request) {
	var users []models.UserInformation
	insertDocs(collectionUsers, writer, request, users)
}

func createMovie(writer http.ResponseWriter, request *http.Request) {
	var movies []models.MovieInformation
	insertDocs(collectionMovies, writer, request, movies)
}
