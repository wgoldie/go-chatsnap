package main

import (
	"encoding/json"
	"fmt"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/kennygrant/sanitize"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/pubnub/go/messaging"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/gopkg.in/dietsche/textbelt.v1"
	"io"
	"net/http"
	"regexp"
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

// Sends log sms to admin (optional)
type TextBeltManager struct {
	Client  *textbelt.Client
	Numbers []string
}

// Function provider to handle http requests to the api for sending messages
func send(bm *BingManager, pn *messaging.Pubnub, tb TextBeltManager) func(w http.ResponseWriter, r *http.Request) {
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

		if cleanQuery == "" || cleanHandle == "" {
			return
		}

		msg := bm.getImageUrls(cleanQuery)

		if err != nil {
			panic(err)
		}

		json, err := json.Marshal(&PubnubMessage{Message: msg, Handle: cleanHandle})
		if err != nil {
			panic(err)
		}

		log := fmt.Sprintf("Chatsnap: %s - %s - %s", m.Channel, m.Handle, m.Channel)
		fmt.Println(log)

		for _, num := range tb.Numbers {
			err := tb.Client.Text(num, log)
			if err != nil {
				fmt.Println("Failed to send text log of message")
				fmt.Println(err)
			}
		}

		var errorChannel = make(chan []byte)
		var callbackChannel = make(chan []byte)
		go pn.Publish(
			m.Channel,
			string(json),
			callbackChannel,
			errorChannel)
	}
}
