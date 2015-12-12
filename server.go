package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"github.com/bndr/gopencils"
	"os"
)

type ImageManager struct {
	Api		*gopencils.Resource	
}

type Keys struct {
	AccountKey string `json:"accountKey"`
}

func (i *ImageManager) getImageUrl(query string) string {
	queryString := map[string]string{
		"AppId": "YOUR_APPID",
		"Version": "2.2",
		"Market": "en-US",
		"Query": "query",
		"Sources": "image",
		"Count": "1"}	
	
	api.Res("

	return queryString["AppId"]
}

type Message struct {
    Message string `json:"message"`
    Handle string `json:"handle"`
    Channel string `json:"channel"`
}

func send(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var m Message
	err := decoder.Decode(&m)
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Println(m.Handle + " " + m.Message + " " + m.Channel)
}



func main() {

	keysFile, _ := os.Open("keys.json")
	decoder := json.NewDecoder(keysFile)
	var k Keys
	err := decoder.Decode(&k)
	if err != nil {
		panic(err)
	}

	auth := gopencils.BasicAuth{k.AccountKey, k.AccountKey}

	im := ImageManager{Api: gopencils.Api("http://api.bing.net/", &auth)}

	fmt.Println("API root:" + im.Api.Url)
	
        http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send)
        http.ListenAndServe(":3000", nil)
}
