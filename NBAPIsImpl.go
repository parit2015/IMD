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
	var targetObj *mongo.Cursor
	targetObj, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Failed to find Documents; ", err)
		return nil, err
	}
	
	return targetObj, nil
}

func findOne(collection *mongo.Collection, filter bson.M, targetObj interface{}) {
	err := collection.FindOne(context.TODO(), filter).Decode(targetObj)
	if err != nil {
		fmt.Println("Failed to find One Document; ", err)
	}
}

func findOneAndUpdate(collection *mongo.Collection, filter bson.M, tobeUpdatedInfo bson.D, targetObj interface{}) {
	err := collection.FindOneAndUpdate(context.TODO(), filter, tobeUpdatedInfo).Decode(targetObj)
	if err != nil {
		fmt.Println("Failed to find/update Document; ", err)
		return
	}
}

func addComment(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var mappingOld, mappingNew models.MovieUserMappingInformation
	var params = mux.Vars(request)
	
	userid, _ := strconv.Atoi(params["user-id"])
	moviename := params["movie-name"]
	findOne(collectionMappings, bson.M{"user": userid, "movie": moviename}, &mappingOld)
	
	_ = json.NewDecoder(request.Body).Decode(&mappingNew)

	for _, v := range mappingNew.Comment {
		mappingOld.Comment = append(mappingOld.Comment, v)
	}
	mappingExtended := mappingOld.Comment

	updateComment := bson.D{
		{"$set", bson.D{
			{"rating", mappingOld.Rating},
			{"comment", mappingExtended},
		}},
	}
	findOneAndUpdate(collectionMappings, 
					bson.M{"user": userid, "movie": moviename}, 
					updateComment, &mappingNew)

	mappingNew.UserId = userid
	mappingNew.MovieName = mux.Vars(request)["movie-name"]
	err := json.NewEncoder(writer).Encode(mappingNew)
	if err != nil {
		return
	}
}

func updateRating(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var mappingOld, mappingNew models.MovieUserMappingInformation
	var params = mux.Vars(request)
	
	userid, _ := strconv.Atoi(params["user-id"])
	moviename := params["movie-name"]
	findOne(collectionMappings, bson.M{"user": userid, "movie": moviename}, &mappingOld)
	
	_ = json.NewDecoder(request.Body).Decode(&mappingNew)
	updateRating := bson.D{
		{"$set", bson.D{
			{"rating", mappingNew.Rating},
			{"comment", mappingOld.Comment},
		}},
	}
	findOneAndUpdate(collectionMappings, bson.M{"user": userid, "movie": moviename}, updateRating, &mappingNew)

	mappingNew.UserId = userid
	mappingNew.MovieName = moviename
	err := json.NewEncoder(writer).Encode(mappingNew)
	if err != nil {
		return
	}
}

func getMoviesByUser(writer http.ResponseWriter, request *http.Request) {
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
		err := mappingsByUserId.Decode(&matchedMapping)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(matchedMapping.MovieName)
		fmt.Println(matchedMapping.Rating)
		
		movieInfo := models.MoviesInfoUserWise{}
		filter := bson.M{"name": matchedMapping.MovieName}
		findOne(collectionMovies, filter, &movieInfo)
		
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
	writer.Header().Set("Content-Type", "application/json")

	var movieInfo models.MovieInformation
	var movieDetailed models.MovieInformationDetailed

	var movieName = mux.Vars(request)["movie-name"]

	movieNameFilter := bson.M{"name": movieName}
	findOne(collectionMovies, movieNameFilter, &movieInfo)

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
