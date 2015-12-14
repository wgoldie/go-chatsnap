package main

import (
	"encoding/json"	
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/kennygrant/sanitize"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/pubnub/go/messaging"
	"io"
	"net/http"
	"regexp"
	"fmt"
)

// Expected format of json for message post requests recieved from clients over the api
type Message struct {
	Message string `json:"message"`
	Handle  string `json:"handle"`
	Channel string `json:"channel"`
}

// Format for json messages broadcast to clients over the pubsub service
type PubnubMessage struct {
	Message []string `json:"images"`
	Handle  string   `json:"sender"`
}

// Function provider to handle http requests to the api for sending messages
func send(im *ImageManager, pn *messaging.Pubnub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m Message
		err := decoder.Decode(&m)
		if err != nil && err != io.EOF {
			panic(err)
		}

		var validChars = regexp.MustCompile(`[^a-zA-Z0-9 ]`)

		sanitizedQuery := sanitize.Accents(m.Message)
		cleanQuery := validChars.ReplaceAllString(sanitizedQuery, "")

		sanitizedHandle := sanitize.Accents(m.Handle)
		cleanHandle := validChars.ReplaceAllString(sanitizedHandle, "")
		
		if cleanQuery == "" || cleanHandle == ""{
			return
		}

		msg := im.getImageUrls(cleanQuery)

		if err != nil {
			panic(err)
		}

		json, err := json.Marshal(&PubnubMessage{Message: msg, Handle: cleanHandle})
		if err != nil {
			panic(err)
		}

		fmt.Println(m.Channel)

		var errorChannel = make(chan []byte)
		var callbackChannel = make(chan []byte)
		go pn.Publish(
			m.Channel,
			string(json),
			callbackChannel,
			errorChannel)
	}
}