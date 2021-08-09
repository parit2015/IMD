package main

import (
	"Sugarbox/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
)

func findMany(collection *mongo.Collection, filter bson.M) (*mongo.Cursor, error) {
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

func findOne(collection *mongo.Collection, filter bson.M, targetObj interface{}) {
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

func getOldNewMappingDocument(request *http.Request, collection *mongo.Collection, 
			DocOld *models.MovieUserMappingInformation,	DocNew *models.MovieUserMappingInformation, filter bson.M) {
	
	findOne(collection, filter, DocOld)
	
	_ = json.NewDecoder(request.Body).Decode(DocNew)
}

func updateDB(collection *mongo.Collection, filter bson.M, tobeUpdateBSONDoc bson.D, 
	targetObj *models.MovieUserMappingInformation) {
	
	findOneAndUpdate(collection, filter, tobeUpdateBSONDoc, targetObj)
}

func addComment(writer http.ResponseWriter, request *http.Request) {
	/*
	This function provides the facility for user to add comment for a particular movie
	
	params:
		userid: passed as the request url variable
		moviename: passed as the request url variable
		comment: passed as the request payload param
	*/
	writer.Header().Set("Content-Type", "application/json")
	
	userid, _ := strconv.Atoi(mux.Vars(request)["user-id"])
	moviename := mux.Vars(request)["movie-name"]

	/*
	1. Old document from db is required to restore other info which is not required to be updated
	2. New document is from request body, it is required in order to fetch the tobe updated param
	*/
	var mappingOld, mappingNew models.MovieUserMappingInformation
	getOldNewMappingDocument(request, collectionMappings, &mappingOld, &mappingNew, 
											bson.M{"user": userid, "movie": moviename})
	
	commentsAppended := append(mappingOld.Comment, mappingNew.Comment...)
	tobeUpdatedBsonDocument := bson.D{
		{"$set", bson.D{
			{"comment", commentsAppended},
		}},
	}
	updateDB(collectionMappings, bson.M{"user": userid, "movie": moviename}, tobeUpdatedBsonDocument, &mappingNew)

	// Update the params to response writer 
	mappingNew.UserId = userid
	mappingNew.MovieName = moviename
	mappingNew.Rating = mappingOld.Rating
	mappingNew.Comment = commentsAppended
	
	err := json.NewEncoder(writer).Encode(mappingNew)
	if err != nil {
		return
	}
}

func updateRating(writer http.ResponseWriter, request *http.Request) {
	/*
	This function provides the facility for user to add/update rating for a particular movie

	description:
		* Movie and User information has been passed in as part of the API url
		* Tobe updated rating information has been passed as the json payload in the request
	*/
	writer.Header().Set("Content-Type", "application/json")

	userid, _ := strconv.Atoi(mux.Vars(request)["user-id"])
	moviename := mux.Vars(request)["movie-name"]
	
	/*
	1. Old document from db is required to restore other info which is not required to be updated
	2. New document is from request body, it is required in order to fetch the tobe updated param
	 */
	var mappingOld, mappingNew models.MovieUserMappingInformation
	getOldNewMappingDocument(request, collectionMappings, &mappingOld, &mappingNew, 
											bson.M{"user": userid, "movie": moviename})
	
	ratingNew := mappingNew.Rating
	tobeUpdatedBsonDocument := bson.D{
		{"$set", bson.D{
			{"rating", ratingNew},
		}},
	}
	updateDB(collectionMappings, bson.M{"user": userid, "movie": moviename}, tobeUpdatedBsonDocument, &mappingNew)

	// Update the params to response writer 
	mappingNew.UserId = userid
	mappingNew.MovieName = moviename
	mappingNew.Comment = mappingOld.Comment
	mappingNew.Rating = ratingNew
	
	err := json.NewEncoder(writer).Encode(mappingNew)
	if err != nil {
		return
	}
}

func getMoviesByUser(writer http.ResponseWriter, request *http.Request) {
	/*
	This function provides the facility for user to see the relevant movie information, in which he/she has either
	added the comment or given a rating

	description:
		* User information has been passed in as part of the API url
	*/
	writer.Header().Set("Content-Type", "application/json")

	moviesByUser := models.MoviesByUserInformation{}
	var params = mux.Vars(request)

	userId, _ := strconv.Atoi(params["user-id"])

	moviesByUser.UserId = userId
	mappingsByUserId, err := findMany(collectionMappings, bson.M{"user": userId})
	if err != nil {
		fmt.Println("Failed to find mapping information for user; ", err)
		return
	}
	for mappingsByUserId.Next(context.TODO()) {
		matchedMapping := models.MovieUserMappingInformation{}
		movieInfo := models.MoviesInfoUserWise{}
		
		err := mappingsByUserId.Decode(&matchedMapping)
		if err != nil {
			log.Fatal(err)
		}
		
		findOne(collectionMovies, bson.M{"name": matchedMapping.MovieName}, &movieInfo)

		movieInfo.Rating = matchedMapping.Rating
		for _, v := range matchedMapping.Comment {
			movieInfo.Comments = append(movieInfo.Comments, v)
		}

		moviesByUser.MoviesInfo = append(moviesByUser.MoviesInfo, movieInfo)
	}
	if err := mappingsByUserId.Err(); err != nil {
		return
	}

	err = json.NewEncoder(writer).Encode(moviesByUser)
	if err != nil {
		return
	}
}

func getMovie(writer http.ResponseWriter, request *http.Request) {
	/*
	This function provides the facility to see the detailed information for all the movies available

	description:
		* Movie information has been passed in as part of the API url
	*/
	writer.Header().Set("Content-Type", "application/json")

	var movieInfo models.MovieInformation
	var movieDetailed models.MovieInformationDetailed

	var movieName = mux.Vars(request)["movie-name"]
	
	findOne(collectionMovies, bson.M{"name": movieName}, &movieInfo)

	movieMappingInfo, err := findMany(collectionMappings, bson.M{"movie": movieInfo.Name})
	if err != nil {
		fmt.Println("Failed to find mapping information for moviename; ", err)
		return
	}
	for movieMappingInfo.Next(context.TODO()) {
		var matchedMapping models.MovieUserMappingInformation
		err := movieMappingInfo.Decode(&matchedMapping)
		if err != nil {
			return
		}

		if matchedMapping.MovieName == movieInfo.Name {
			movieDetailed.MovieInfo = models.MovieInformation{Name: movieInfo.Name, Type: movieInfo.Type, Description: movieInfo.Description}
			movieDetailed.Count += 1
			movieDetailed.Rating = (movieDetailed.Rating + matchedMapping.Rating) / float32(movieDetailed.Count)
			for _, v := range matchedMapping.Comment {
				movieDetailed.Comments = append(movieDetailed.Comments, v)
			}
		}
	}
	if err := movieMappingInfo.Err(); err != nil {
		return
	}

	err = json.NewEncoder(writer).Encode(movieDetailed)
	if err != nil {
		return
	}
}
