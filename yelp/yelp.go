package yelp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	mClientId = "rTo1ltY1Ct63KEv-SjR2vg"
	mClientSecret = "uUwVcqkDSXEhJE8sCHBHc6JXUBivt40BxSde5Mmdlmzudloov5F6UKPWTBICgDbI"
	mUrlAccessToken = "https://api.yelp.com/oauth2/token"
	mUrlBusinessSearch = "https://api.yelp.com/v3/businesses/search"
)

type ParamSearchBusiness struct
{
	Term string `json:"term"`
	Location string `json:"location"`
}

type RespSearchBusiness struct {
	Term 		string `json:"term"`
	Location 	string `json:"location"`
	Total      	int64  `json:"total"`
	Businesses 	[]Business `json:"businesses"`
	Region 		struct {
				Center struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
				} `json:"center"`
			} `json:"region"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type Business struct {
	Rating 		float32 `json:"rating"`
	Price 		string  `json:"price"`
	Id 		string  `json:"id"`
	IsClosed 	bool `json:"is_closed"`
	Categories 	[]struct {
					Alias string `json:"alias"`
					Title string `json:"title"`
	} `json:"categories"`
	ReviewCount	int32 `json:"review_count"`
	Name 		string  `json:"name"`
	Url 		string  `json:"url"`
	Coordinates 	struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	ImgUrl 		string  `json:"image_url"`
	Location	struct {
					City 		string `json:"city"`
					Country 	string `json:"country"`
					Address2 	string `json:"address2"`
					Address3 	string `json:"address3"`
					State 		string `json:"state"`
					Address1 	string `json:"address1"`
					ZipCode 	string `json:"zip_code"`
	} `json:"location"`
	Distance	float32 `json:"distance"`
	Transactions	[]string `json:"transactions"`
}

//Get Yelp access token
func GetAccessToken(pAccessToken *AccessToken) error {
	lUrlValues := url.Values{}
	lUrlValues.Add("grant_type", "client_credentials")
	lUrlValues.Add("client_id", mClientId)
	lUrlValues.Add("client_secret", mClientSecret)

	lBufferUrlValues := bytes.NewBuffer([]byte(lUrlValues.Encode()))

	//POST request for getting AccessToken. The respose will Access token with its type
	lResponse, lErr := http.Post(mUrlAccessToken, "application/x-www-form-urlencoded", lBufferUrlValues)
	if lErr != nil {
		return lErr
	}

	defer lResponse.Body.Close()
	lByteArr, lErr := ioutil.ReadAll(lResponse.Body)
	if lErr != nil {
		return lErr
	}

	lErr = json.Unmarshal(lByteArr, &pAccessToken)
	if lErr != nil {
		return lErr
	}

	return nil
}

//Yelp business search
func SearchBusiness(pParamSearchBusiness ParamSearchBusiness, pAccessTokenPtr *AccessToken) (*RespSearchBusiness, error) { 
	lUrlWithParam := mUrlBusinessSearch + "?" + "term=" + pParamSearchBusiness.Term + "&location=" + pParamSearchBusiness.Location

	lRequest, lErr := http.NewRequest("GET", lUrlWithParam, nil)
	if lErr != nil {
		return nil, lErr
	}

	lRequest.Header.Set("Authorization", pAccessTokenPtr.TokenType + " " + pAccessTokenPtr.AccessToken)

	//GET request for retrieving Restaurant. The results will be with all details
	lRefHttpClient := &http.Client{}
	lResponse, lErr := lRefHttpClient.Do(lRequest)
	if lErr != nil {
		return nil, lErr
	}

	lByteArr, lErr := ioutil.ReadAll(lResponse.Body)
	if lErr != nil {
		return nil, lErr
	}

	var lRespSearchBusiness = new(RespSearchBusiness)
	lErr = json.Unmarshal([]byte(lByteArr), &lRespSearchBusiness)
	if lErr != nil {
		return nil, lErr
	}

	//to make the response with term and location 
	lRespSearchBusiness.Term = pParamSearchBusiness.Term
	lRespSearchBusiness.Location = pParamSearchBusiness.Location

	return lRespSearchBusiness, nil
}
