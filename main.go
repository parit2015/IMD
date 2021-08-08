package main

import (
	"Sugarbox/helper"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var collectionMovies, collectionUsers, 
	collectionMappings = helper.GetCollection(helper.GetDBConnection("mongodb://localhost:27017"))

func main() {
	r := mux.NewRouter()
	
	/*
	These endpoints are used to populate/access information to/from DB. 
	As per requirement, we dont want this functionality to exposed from the Application.
	Keeping it internal.
	 */
	r.HandleFunc("/IMDB/movies", createMovie).Methods("POST")
	r.HandleFunc("/IMDB/users", createUsers).Methods("POST")
	r.HandleFunc("/IMDB/mappings", createMappings).Methods("POST")
	//r.HandleFunc("/IMDB/deleteCollections", deleteCollectionMovies).Methods("DELETE")
	//r.HandleFunc("/IMDB/getMappings", getMappings).Methods("GET")
	//r.HandleFunc("/IMDB/getUsers", getUsers).Methods("GET")
	//r.HandleFunc("/IMDB/getMovies", getMovies).Methods("GET")
	
	/*
	The set of endpoints to be exposed as northbound APIs
	 */
	r.HandleFunc("/IMDB/movies/{movie-name}", getMovie).Methods("GET")
	r.HandleFunc("/IMDB/moviesByUser/{user-id}", getMoviesByUser).Methods("GET")
	r.HandleFunc("/IMDB/updateRating/{user-id}/{movie-name}", updateRating).Methods("PUT")
	r.HandleFunc("/IMDB/addComment/{user-id}/{movie-name}", addComment).Methods("PUT")
	
	log.Fatal(http.ListenAndServe(":8000", r))
}