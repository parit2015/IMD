package utils

import (
	"IMD-master/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func GetOldNewDocument(request *http.Request, collection *mongo.Collection,
	DocOld *models.MovieUserMappingInformation, DocNew *models.MovieUserMappingInformation, filter bson.M) {

	FindOne(collection, filter, DocOld)

	_ = json.NewDecoder(request.Body).Decode(DocNew)
}

func UpdateDB(collection *mongo.Collection, filter bson.M, tobeUpdateBSONDoc bson.D,
	targetObj *models.MovieUserMappingInformation) {

	findOneAndUpdate(collection, filter, tobeUpdateBSONDoc, targetObj)
}
