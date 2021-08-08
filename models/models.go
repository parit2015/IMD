package models

type MovieInformation struct {
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type UserInformation struct {
	Id    int    `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Addr  string `json:"addr" bson:"addr"`
}

type MovieUserMappingInformation struct {
	UserId    	int      `json:"user" bson:"user"`
	MovieName   string   `json:"movie" bson:"movie"`
	Rating  	float32   `json:"rating,omitempty" bson:"rating,omitempty"`
	Comment 	[]string `json:"comment,omitempty" bson:"comment,omitempty"`
}

type MovieInformationDetailed struct {
	MovieInfo	MovieInformation
	Rating  	float32
	Count		int
	Comments	[]string
}

type MoviesByUserInformation struct {
	UserId 		int
	MoviesInfo 	[]MoviesInfoUserWise
}

type MoviesInfoUserWise struct {
	Name		string
	Rating  	float32
	Comments	[]string
}