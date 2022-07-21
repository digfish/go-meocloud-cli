package meocloudcli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"encoding/json"
	"net/http"

	"net/url"

	"github.com/mrjones/oauth"
)

// Language: go
// Path: main.go
// Returns an http.Client ready to be used to interact with the MeoCloud API
// Requires an .env file with the following variables:
// CONSUMER_KEY CONSUMER_SECRET OAUTH_TOKEN OAUTH_SECRET
// This file should be placed in the main dir
// Note: at this time, the process of obtaining an access token is not yet implemented
func get_meo_client() *http.Client {
	godotenv.Load()
	accessToken := &oauth.AccessToken{
		Token:  os.Getenv("OAUTH_TOKEN"),
		Secret: os.Getenv("OAUTH_TOKEN_SECRET"),
	}
	consumer := oauth.NewConsumer(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://meocloud.pt/oauth/request_token",
			AuthorizeTokenUrl: "https://meocloud.pt/oauth/authorize",
			AccessTokenUrl:    "https://meocloud.pt/oauth/access_token"},
	)
	meo, _ := consumer.MakeHttpClient(accessToken)
	return meo
}

// Given an url, returns a map of the json response
// To be used by API calls that return JSON
func fetch_json(url string) (map[string]interface{}, int) {
	bodybytes, status := fetch(url)
	var bodymap map[string]interface{}
	json.Unmarshal(bodybytes, &bodymap)
	return bodymap, status
}

// Generic function to fetch a url and return the body as a byte array
func fetch(url string) ([]byte, int) {
	meo := get_meo_client()
	httpResponse, _ := meo.Get(url)
	bodybytes, _ := io.ReadAll(httpResponse.Body)
	return bodybytes, httpResponse.StatusCode
}

// Get the account information of the current user (https://meocloud.pt/documentation#accountinfo)
func account_info() (map[string]interface{}, int) {
	return fetch_json("https://api.meocloud.pt/1/Account/Info")
}

// for debugging purposes, outputs the content of a map obtained from a JSON response
func dump_json(bodymap map[string]interface{}) {
	for key, value := range bodymap {
		fmt.Println(key, "->", value)
	}
}

// get the metadata of a file or folder, used of directory listings (https://meocloud.pt/documentation#metadata)
func get_metadata(path string) (map[string]interface{}, int) {
	bodymap, status := fetch_json(fmt.Sprintf("https://api.meocloud.pt/1/Metadata/meocloud%s", path))

	/* 	for key,value := range bodymap["contents"].([]interface{})  {
	   		fmt.Println(key,"->",value.(map[string]interface{})["name"])
	   	}
	*/
	return bodymap, status
}

// Get the contents of a file given its path in the cloud (https://meocloud.pt/documentation#files)
func get_file(filepath string) ([]byte, int) {
	meo := get_meo_client()
	httpResponse, _ := meo.Get(fmt.Sprintf("https://api-content.meocloud.pt/1/Files/meocloud/%s", filepath))
	bodybytes, _ := io.ReadAll(httpResponse.Body)
	return bodybytes, httpResponse.StatusCode
}

// Send byte array to the cloud (https://meocloud.pt/documentation#files) giving it a name specified in newfilepath
// returns the status code
func send(newfilepath string, data []byte) int {
	meo := get_meo_client()
	putReq, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://api-content.meocloud.pt/1/Files/meocloud/%s", newfilepath),
		bytes.NewReader(data))
	httpResponse, _ := meo.Do(putReq)
	return httpResponse.StatusCode
}

// send a file to the cloud (https://meocloud.pt/documentation#files)
func send_file(filepath string) int {
	contentBytes, _ := os.ReadFile(filepath)
	return send(filepath, contentBytes)
}

// delete a file in the cloud (https://meocloud.pt/documentation#delete)
// returns the status code (200 on success, 404 if the file does not exist and 406 if many files are deleted)
func delete_file(filepath string) int {
	meo := get_meo_client()
	httpResponse, _ := meo.PostForm("https://api.meocloud.pt/1/Fileops/Delete",
		url.Values{
			"root": {"meocloud"},
			"path": {strings.Join([]string{"/", filepath}, "")},
		})
	return httpResponse.StatusCode
}

// creates a directory in the cloud (https://meocloud.pt/documentation#createfolder)
// returns the status code (200 if successful, 403 if the directory already exists)
func create_dir(folderpath string) int {
	meo := get_meo_client()
	httpResponse, _ := meo.PostForm("https://api.meocloud.pt/1/Fileops/CreateFolder",
		url.Values{
			"root": {"meocloud"},
			"path": {strings.Join([]string{"/", folderpath}, "")},
		})
	return httpResponse.StatusCode

}
