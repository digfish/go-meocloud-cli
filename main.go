package main

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

func fetch_json(url string) (map[string]interface{},int) {
	bodybytes,status := fetch(url)
	var bodymap map[string]interface{}
	json.Unmarshal(bodybytes,&bodymap)
	return bodymap,status;
}

func fetch(url string) ([]byte,int)  {
	meo := get_meo_client()
	httpResponse, _ := meo.Get(url)
	bodybytes,_ := io.ReadAll(httpResponse.Body)
	return bodybytes, httpResponse.StatusCode
}

func account_info() (map[string]interface{},int) {
	return  fetch_json("https://api.meocloud.pt/1/Account/Info")
}

func dump_json(bodymap map[string]interface{}) {
		for key,value := range bodymap {
		fmt.Println(key,"->",value)
	}
}

func get_metadata(path string)  (map[string]interface{},int) {
	bodymap,status := fetch_json( fmt.Sprintf("https://api.meocloud.pt/1/Metadata/meocloud%s",path))

/* 	for key,value := range bodymap["contents"].([]interface{})  {
		fmt.Println(key,"->",value.(map[string]interface{})["name"])
	}
 */
	return bodymap,status
}

func get_file(filepath string) ([]byte,int) {
	meo := get_meo_client()
	httpResponse, _ := meo.Get(fmt.Sprintf("https://api-content.meocloud.pt/1/Files/meocloud/%s",filepath))
	bodybytes,_ := io.ReadAll(httpResponse.Body)
	return bodybytes, httpResponse.StatusCode
}

func send(newfilepath string,data []byte) (int) {
	meo := get_meo_client()
	putReq,_ := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://api-content.meocloud.pt/1/Files/meocloud/%s",newfilepath),
		bytes.NewReader(data))
	httpResponse, _ := meo.Do(putReq)
	return httpResponse.StatusCode
}

func send_file(filepath string) (int) {
	contentBytes,_ := os.ReadFile(filepath)
	return send(filepath, contentBytes)
}

func delete_file(filepath string) (int) {
	meo := get_meo_client()
	httpResponse, _ := meo.PostForm("https://api.meocloud.pt/1/Fileops/Delete",
		url.Values{
			"root": {"meocloud"},
			"path": {strings.Join([]string{"/",filepath},"")},
		})
	return httpResponse.StatusCode
}

func create_dir(folderpath string) (int) {
		meo := get_meo_client()
	httpResponse, _ := meo.PostForm("https://api.meocloud.pt/1/Fileops/CreateFolder",
		url.Values{
			"root": {"meocloud"},
			"path": {strings.Join([]string{"/",folderpath},"")},
		})
	return httpResponse.StatusCode

}

func main() {
	md,status :=get_metadata("/")
	dump_json(md)
	fmt.Println(status)
}
