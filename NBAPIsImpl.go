package main

import (
	"IMD-master/models"
	"IMD-master/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
)

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
	utils.GetOldNewDocument(request, collectionMappings, &mappingOld, &mappingNew, 
											bson.M{"user": userid, "movie": moviename})
	
	commentsAppended := append(mappingOld.Comment, mappingNew.Comment...)
	tobeUpdatedBsonDocument := bson.D{
		{"$set", bson.D{
			{"comment", commentsAppended},
		}},
	}
	utils.UpdateDB(collectionMappings, bson.M{"user": userid, "movie": moviename}, tobeUpdatedBsonDocument, &mappingNew)

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
	utils.GetOldNewDocument(request, collectionMappings, &mappingOld, &mappingNew, 
											bson.M{"user": userid, "movie": moviename})
	
	ratingNew := mappingNew.Rating
	tobeUpdatedBsonDocument := bson.D{
		{"$set", bson.D{
			{"rating", ratingNew},
		}},
	}
	utils.UpdateDB(collectionMappings, bson.M{"user": userid, "movie": moviename}, tobeUpdatedBsonDocument, &mappingNew)

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
	mappingsByUserId, err := utils.FindMany(collectionMappings, bson.M{"user": userId})
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
		
		utils.FindOne(collectionMovies, bson.M{"name": matchedMapping.MovieName}, &movieInfo)

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
	
	utils.FindOne(collectionMovies, bson.M{"name": movieName}, &movieInfo)

	movieMappingInfo, err := utils.FindMany(collectionMappings, bson.M{"movie": movieInfo.Name})
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
