package main

import (
	"Sugarbox/models"
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
