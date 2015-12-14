package main

import (
	"encoding/json"
	"fmt"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/pubnub/go/messaging"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Message string `json:"message"`
	Handle  string `json:"handle"`
	Channel string `json:"channel"`
}

type PubnubMessage struct {
	Message []string `json:"images"`
	Handle  string   `json:"sender"`
}

func send(im *ImageManager, pn *messaging.Pubnub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m Message
		err := decoder.Decode(&m)
		if err != nil && err != io.EOF {
			panic(err)
		}

		msg := im.getImageUrl(m.Message)

		json, err := json.Marshal(&PubnubMessage{Message: []string{msg}, Handle: m.Handle})
		if err != nil {
			panic(err)
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

func main() {
	yahooClientId := os.Getenv("YAHOO_CLIENT_ID")
	yahooClientSecret := os.Getenv("YAHOO_CLIENT_SECRET")
	pubnubPublishKey := os.Getenv("PUBNUB_PUBLISH_KEY")
	pubnubSubscribeKey :=  os.Getenv("PUBNUB_SUBSCRIBE_KEY")
	pubnubSecretKey :=  os.Getenv("PUBNUB_SECRET_KEY")
	port := os.Getenv("PORT")
	if yahooClientId == "" || yahooClientSecret == "" || pubnubPublishKey == "" || pubnubSubscribeKey == "" || pubnubSecretKey == "" || port == "" {
		fmt.Println("Something is wrong with the config flags")
		os.Exit(666)
	}
	

	pn := messaging.NewPubnub(pubnubPublishKey, pubnubSubscribeKey, pubnubSecretKey, "", false, "92895fc3-cc14-4e3d-a38a-901dd3739238")

	im := NewImageManager(yahooClientId, yahooClientSecret)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send(&im, pn))
	http.ListenAndServe(":"+port, nil)
}
