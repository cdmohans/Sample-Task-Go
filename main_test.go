package main

import (
	"testing" 
	"./yelp"
	"./mongo"
	"log"
)

var (
	tYelpAccessToken *yelp.AccessToken
	tYelpParamSearchBusiness yelp.ParamSearchBusiness
	tYelpRespSearchBusiness yelp.RespSearchBusiness
)


//Validating the Yelp API AccessToken
func TestAccessToken(t *testing.T) {
	tYelpAccessToken = new(yelp.AccessToken)
	lErr := yelp.GetAccessToken(tYelpAccessToken)
	if lErr != nil {
		t.Errorf("Error retrieving AccessToken")
	}
}

//Validating the Yelp API searchbusiness call
func TestSearchBusiness(t *testing.T) {
	tYelpParamSearchBusiness.Term = "coffee"
	tYelpParamSearchBusiness.Location = "singapore"
	tYelpRespSearchBusiness, lErr := yelp.SearchBusiness(tYelpParamSearchBusiness, tYelpAccessToken)
	
	if lErr != nil {
		t.Errorf("Error retrieving AccessToken")
	}
	
	if tYelpRespSearchBusiness.Total == 0 {
		t.Errorf("Error on Response of the YelpSearchParameters coffee and singapore")
	}
}


//Validating the Yelp SearchBusiness call with wrong input
func TestSearchBusiness_invalid(t *testing.T) {
	var lYelpParamSearchBusiness yelp.ParamSearchBusiness
	lYelpParamSearchBusiness.Term = "ghfd"
	lYelpParamSearchBusiness.Location = "ijkl" // junk values 
	lYelpRespSearchBusiness, lErr := yelp.SearchBusiness(lYelpParamSearchBusiness, tYelpAccessToken)
	
	if lErr != nil {
		t.Errorf("Error retrieving AccessToken")
	}
	
	if lYelpRespSearchBusiness.Total != 0 {
		t.Errorf("Error on Response of the YelpSearchParameters coffee and singapore")
	}
}


//Testing MongoDB insertion and retrieving
func TestMongoDb(t *testing.T) {
	lInsertIntoMongo, lErr := mongo.CreateRespSearchBusiness(&tYelpRespSearchBusiness)
	if lErr != nil {
	    t.Errorf("Error on inserting into Mongo Db")
	}
	log.Println("Insert into Mongo: ", lInsertIntoMongo)

	lYelpRespSearchBusinessFromMongoDB, lErr := mongo.GetRespSearchBusiness(tYelpParamSearchBusiness.Term, tYelpParamSearchBusiness.Location)
	if lErr != nil {
	    t.Errorf("Error on retrieving from Mongo Db")
	}
	log.Println("Retrieved from Mongo: ", lYelpRespSearchBusinessFromMongoDB)
	
	if lErr == nil {
		if lYelpRespSearchBusinessFromMongoDB.Term != tYelpParamSearchBusiness.Term {
		    t.Errorf("Retrieving and inserted records are no same.`")
			return
		}
	}

}