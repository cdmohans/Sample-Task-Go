package mongo

import (
    "log"
	"../yelp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Insert business details into mongo
func CreateRespSearchBusiness(pRespSearchBusiness  *yelp.RespSearchBusiness) (bool, error) {
	lSession, lErr := mgo.Dial("localhost:27017")
	if lErr != nil {
		log.Println("mgo.Dial is failed")
		return false, lErr
	}
	defer lSession.Close()

	//table will be RespSearchBusiness
	lCollection := lSession.DB("yelp").C("RespSearchBusiness")
	lErr = lCollection.Insert(pRespSearchBusiness)
	if lErr != nil {
		log.Println("Insert in mongo is failed.")
		return false, lErr
	}
	return true, lErr
}

//Retrieving business details from mongo
func GetRespSearchBusiness(pTerm string, pLocation string) (*yelp.RespSearchBusiness, error) {
	var lRespSearchBusiness = new(yelp.RespSearchBusiness)

	lSession, lErr := mgo.Dial("localhost:27017")
	if lErr != nil {
		log.Println("mgo.Dial is failed")
		return lRespSearchBusiness, lErr
	}
	defer lSession.Close()	

	lCollection := lSession.DB("yelp").C("RespSearchBusiness")
	lErr = lCollection.Find(bson.M{"term": pTerm, "location": pLocation}).One(&lRespSearchBusiness)
	if lErr != nil {
		log.Println("Find in mongo is failed.")
		return lRespSearchBusiness, lErr
	}
	return lRespSearchBusiness, lErr
}
