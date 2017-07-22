package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"./yelp"
	"./mongo"
	//"testing"
)

const (
	mFolderLibs string = "libs"
	mFolderHtml string = "html"
)

var (
	mYelpAccessTokenTimeExpired time.Time 
	mYelpAccessToken *yelp.AccessToken 
)

//Render Html pages
func RenderPage(pResponseWriter http.ResponseWriter, pHtmlFileName string) {
	log.Println("Processing " + pHtmlFileName + ".html... " )
	lListHtmlFile := []string{ mFolderHtml + "/base.html", 
				mFolderHtml + "/" + pHtmlFileName + ".html" }

	lTemplatePtr, lErr := template.ParseFiles(lListHtmlFile...)
	if lErr != nil {
		log.Println("Error parsing files: ", lErr)
	}

	lErr = lTemplatePtr.Execute(pResponseWriter, nil)
	if lErr != nil {
		log.Println("Error executing template: ", lErr)
	}
}

//Handler for processing files in libs folder
func HandlerFilesLib(pResponseWriter http.ResponseWriter, pHttpRequest *http.Request) {
	lPathLibFile := pHttpRequest.URL.Path[len("/" + mFolderLibs + "/"):]
	
	if len(lPathLibFile) != 0 {
		lFile, lErr := http.Dir(mFolderLibs + "/").Open(lPathLibFile)
		if lErr != nil {
			log.Println("Error opening lib file: ", lErr)
			http.NotFound(pResponseWriter, pHttpRequest)
		} else {
			lReadSeeker := io.ReadSeeker(lFile)
			http.ServeContent(pResponseWriter, pHttpRequest, lPathLibFile, time.Now(), lReadSeeker)
		}
	}
}

//Handler for processing html/index.html
func HandlerPageIndex(pResponseWriter http.ResponseWriter, pHttpRequest *http.Request) {
	RenderPage(pResponseWriter, "index")
}

//Handler for searching business in yelp
func HandlerSearchBusinessInYelp(pResponseWriter http.ResponseWriter, pRequest *http.Request) {
	pResponseWriter.Header().Set("Content-Type", "application/json")

	var lYelpParamSearchBusiness yelp.ParamSearchBusiness 
	lDecoder := json.NewDecoder(pRequest.Body)
	lErr := lDecoder.Decode(&lYelpParamSearchBusiness)
	if lErr != nil {
	        log.Println("Error decoding Json string: ", lErr)
	}
	defer pRequest.Body.Close()


	//Retrieving info of term and location from mongo. This will minimize the yelp calling
	lYelpRespSearchBusinessFromMongoDB, lErr := mongo.GetRespSearchBusiness(lYelpParamSearchBusiness.Term, lYelpParamSearchBusiness.Location)
	if lErr == nil {
		if lYelpRespSearchBusinessFromMongoDB.Term == lYelpParamSearchBusiness.Term {
			json.NewEncoder(pResponseWriter).Encode(lYelpRespSearchBusinessFromMongoDB)
		        log.Println("Retrieved data from Mongo DB...")
			return
		}
	}

	//Avoid using Expired AccessToken by calling SetYelpAccessToken if it is expired
	if IsYelpAccessTokenExpired() {
		SetYelpAccessToken()
	}

	log.Println("Yelp access token expiring on " + mYelpAccessTokenTimeExpired.String() + "...") 

	//If mongo does not have required info of term and location, search it by Yelp
	lYelpRespSearchBusiness, lErr := yelp.SearchBusiness(lYelpParamSearchBusiness, mYelpAccessToken)
    	if lErr != nil {
	        log.Println("Error searching businesses in Yelp: ", lErr)
    	}

	json.NewEncoder(pResponseWriter).Encode(lYelpRespSearchBusiness)

	//Insert it in mongo as it is not found already in mongo
	lInsertIntoMongo, lErr := mongo.CreateRespSearchBusiness(lYelpRespSearchBusiness)
	if lErr != nil {
	        log.Println("Error decoding Json string: ", lErr)
	}

	log.Println("Insert into Mongo: ", lInsertIntoMongo)
}

//Returns true if the Access Token of Yelp is expired.
func IsYelpAccessTokenExpired() (bool) {
	if time.Now().After(mYelpAccessTokenTimeExpired) {
		log.Println("Yelp access token expired on " + mYelpAccessTokenTimeExpired.String()) 
		return true
	}
	return false
}

//Retrieving the access token for yelp.
func SetYelpAccessToken() {
    log.Println("Retrieving Yelp access token...")
	mYelpAccessTokenTimeExpired = time.Now()
	mYelpAccessToken = new(yelp.AccessToken)
	lErr := yelp.GetAccessToken(mYelpAccessToken)
	if lErr != nil {
		log.Println("Error retrieving Yelp access token: ", lErr)
	}
	lDuration := time.Duration(mYelpAccessToken.ExpiresIn*1000*1000*1000)
	mYelpAccessTokenTimeExpired = mYelpAccessTokenTimeExpired.Add(lDuration)
}

//Main function
func main(){
	//Get the access token for yelp
	SetYelpAccessToken()

	//Handler functions registration.
	//SearchBusiness will be called from index.html
    log.Println("Registering handlers...")
	http.HandleFunc("/" + mFolderLibs + "/", HandlerFilesLib)
	http.HandleFunc("/", HandlerPageIndex)
	http.HandleFunc("/SearchBusiness", HandlerSearchBusinessInYelp)
	
	//Localhost:8080 listening
    log.Println("Starting webserver...")
	lErrListenAndServe := http.ListenAndServe(":8000", nil)
	if lErrListenAndServe != nil {
	        log.Println("Error starting webserver: ", lErrListenAndServe)
	}
}

