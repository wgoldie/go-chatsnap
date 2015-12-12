package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"os"
	"flag"
)

type Message struct {
    Message string `json:"message"`
    Handle string `json:"handle"`
    Channel string `json:"channel"`
}

func send(im *ImageManager) func (w http.ResponseWriter, r *http.Request)  {
	return func (w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m Message
		err := decoder.Decode(&m)
		if err != nil && err != io.EOF {
			panic(err)
		}
	
		fmt.Println(m.Handle + " " + m.Message + " " + m.Channel)
		fmt.Println("BUG: " + im.getImageUrl(m.Message))
	}
}



func main() {
	clientId := flag.String("yahooClientId", "", "The yahoo client ID key")
	clientSecret := flag.String("yahooClientSecret", "", "The yahoo client secret key")
	
	flag.Parse()
	
	if *clientId == "" || *clientSecret == "" {
		fmt.Println(*clientId)
		fmt.Println(*clientSecret)
		os.Exit(1)
	}
	
	im := NewImageManager(*clientId, *clientSecret)

	fmt.Println("API root:" + im.Url)
        http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send(&im))
        http.ListenAndServe(":3000", nil)
}
