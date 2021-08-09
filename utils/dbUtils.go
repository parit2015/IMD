package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func FindMany(collection *mongo.Collection, filter bson.M) (*mongo.Cursor, error) {
	/*
	This function finds more than one document

	params:
		collection: Collection on which the search has to be performed
		filter: Filtering condition of the search

	returns:
		mongo-cursor: Collection of resultant document
	*/
	var targetObj *mongo.Cursor
	targetObj, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Failed to find Documents; ", err)
		return nil, err
	}

	return targetObj, nil
}

func FindOne(collection *mongo.Collection, filter bson.M, targetObj interface{}) {
	/*
	This function finds one and only one document

	params:
		collection: Collection on which the search has to be performed
		filter: Filtering condition of the search
		targetObj: Object on which the resultant document has to be saved
	*/
	err := collection.FindOne(context.TODO(), filter).Decode(targetObj)
	if err != nil {
		fmt.Println("Failed to find One Document; ", err)
	}
}

func findOneAndUpdate(collection *mongo.Collection, filter bson.M, tobeUpdatedInfo bson.D, targetObj interface{}) {
	/*
	This function finds and update one document

	params:
		collection: Collection on which the search has to be performed
		filter: Filtering condition of the search
		tobeUpdatedInfo: The payload, that needs to be updated in the filtered document
	*/
	err := collection.FindOneAndUpdate(context.TODO(), filter, tobeUpdatedInfo).Decode(targetObj)
	if err != nil {
		fmt.Println("Failed to find/update Document; ", err)
		return
	}
}

func InsertDocs(collection *mongo.Collection, writer http.ResponseWriter, request *http.Request,
					targetObj interface{}) {
	writer.Header().Set("Content-Type", "application/json")

	_ = json.NewDecoder(request.Body).Decode(&targetObj)
	for _, v := range targetObj.([]interface{}) {
		result, err := collection.InsertOne(context.TODO(), v)
		if err != nil {
			fmt.Println("Failed to insert information; ", err)
			return
		}

		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			return
		}
	}
}
