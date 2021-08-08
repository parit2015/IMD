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

func addComment(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(request)
	userid, _ := strconv.Atoi(params["user-id"])
	moviename := params["movie-name"]

	filter := bson.M{"user": userid, "movie": moviename}

	var mappingOld models.MovieUserMappingInformation
	err := collectionMappings.FindOne(context.TODO(), filter).Decode(&mappingOld)
	if err != nil {
		fmt.Println("Failed to find Mapping information for userid: ", userid, " and moviename: ", moviename, 
			"; ", err)
		return
	}

	var mapping models.MovieUserMappingInformation
	_ = json.NewDecoder(request.Body).Decode(&mapping)

	for _, v := range mapping.Comment {
		mappingOld.Comment = append(mappingOld.Comment, v)
	}
	mappingExtended := mappingOld.Comment

	update := bson.D{
		{"$set", bson.D{
			{"rating", mappingOld.Rating},
			{"comment", mappingExtended},
		}},
	}

	err = collectionMappings.FindOneAndUpdate(context.TODO(), filter, update).Decode(&mapping)
	if err != nil {
		fmt.Println("Failed to find/update Mapping information for userid: ", userid, " and moviename: ", moviename,
			"; ", err)
		return
	}

	mapping.UserId = userid
	mapping.MovieName = moviename
	err = json.NewEncoder(writer).Encode(mapping)
	if err != nil {
		return
	}
}

func updateRating(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(request)
	userid, _ := strconv.Atoi(params["user-id"])
	moviename := params["movie-name"]

	filter := bson.M{"user": userid, "movie": moviename}

	var mappingOld models.MovieUserMappingInformation
	err := collectionMappings.FindOne(context.TODO(), filter).Decode(&mappingOld)
	if err != nil {
		fmt.Println("Failed to find Mapping information for userid: ", userid, " and moviename: ", moviename,
			"; ", err)
		return
	}

	var mapping models.MovieUserMappingInformation
	_ = json.NewDecoder(request.Body).Decode(&mapping)

	update := bson.D{
		{"$set", bson.D{
			{"rating", mapping.Rating},
			{"comment", mappingOld.Comment},
		}},
	}

	err = collectionMappings.FindOneAndUpdate(context.TODO(), filter, update).Decode(&mapping)
	if err != nil {
		fmt.Println("Failed to find/update Mapping information for userid: ", userid, " and moviename: ", moviename,
			"; ", err)
		return
	}

	mapping.UserId = userid
	mapping.MovieName = moviename
	err = json.NewEncoder(writer).Encode(mapping)
	if err != nil {
		return
	}
}

func getMoviesByUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	moviesByUser := models.MoviesByUserInformation{}

	var userId, _ = strconv.Atoi(mux.Vars(request)["user-id"])
	moviesByUser.UserId = userId
	cur, err := collectionMappings.Find(context.TODO(), bson.M{"user": userId})
	if err != nil {
		fmt.Println("Failed to find Mapping information for userid: ", userId, "; ", err)
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.TODO())

	for cur.Next(context.TODO()) {
		matchedMapping := models.MovieUserMappingInformation{}
		err := cur.Decode(&matchedMapping)
		if err != nil {
			log.Fatal(err)
		}

		movieInfo := models.MoviesInfoUserWise{}
		err = collectionMovies.FindOne(context.TODO(), bson.M{"name": matchedMapping.MovieName}).Decode(&movieInfo)
		if err != nil {
			fmt.Println("Failed to find Movie information for moviename: ", matchedMapping.MovieName, "; ", err)
			return
		}
		movieInfo.Rating = matchedMapping.Rating
		for _, v := range matchedMapping.Comment {
			movieInfo.Comments = append(movieInfo.Comments, v)
		}

		moviesByUser.MoviesInfo = append(moviesByUser.MoviesInfo, movieInfo)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(writer).Encode(moviesByUser)
	if err != nil {
		return
	}
}

func getMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var movieInfo models.MovieInformation
	
	var movieName = mux.Vars(request)["movie-name"]

	movieNameFilter := bson.M{"name": movieName}
	err := collectionMovies.FindOne(context.TODO(), movieNameFilter).Decode(&movieInfo)
	if err != nil {
		fmt.Println("Failed to find Movie information for moviename: ", movieName, "; ", err)
		return
	}
	
	movieMappingInfo, err := collectionMappings.Find(context.TODO(), bson.M{"movie": movieInfo.Name})
	if err != nil {
		fmt.Println("Failed to find Mapping information for moviename: ", movieInfo.Name, "; ", err)
		return
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(movieMappingInfo, context.TODO())

	var movieDetailed models.MovieInformationDetailed
	for movieMappingInfo.Next(context.TODO()) {
		var matchedMapping models.MovieUserMappingInformation
		err := movieMappingInfo.Decode(&matchedMapping)
		if err != nil {
			log.Fatal(err)
		}

		if matchedMapping.MovieName == movieInfo.Name {
			movieDetailed.MovieInfo = models.MovieInformation{Name: movieInfo.Name, Type: movieInfo.Type, Description: movieInfo.Description}
			movieDetailed.Count += 1
			movieDetailed.Rating = (movieDetailed.Rating + matchedMapping.Rating)/float32(movieDetailed.Count)
			for _, v := range matchedMapping.Comment {
				movieDetailed.Comments = append(movieDetailed.Comments, v)
			}
		}
	}
	if err := movieMappingInfo.Err(); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(writer).Encode(movieDetailed)
	if err != nil {
		return
	}
}
