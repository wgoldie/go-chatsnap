package main

import (
	"fmt"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/pubnub/go/messaging"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/satori/go.uuid"
	"net/http"
	"os"
)

func main() {
	yahooClientId := os.Getenv("YAHOOBOSS_KEY")
	yahooClientSecret := os.Getenv("YAHOOBOSS_SECRET")
	pubnubPublishKey := os.Getenv("PUBNUB_PUBLISH_KEY")
	pubnubSubscribeKey := os.Getenv("PUBNUB_SUBSCRIBE_KEY")
	pubnubSecretKey := os.Getenv("PUBNUB_SECRET_KEY")
	redisUrl := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")

	if yahooClientId == "" || yahooClientSecret == "" || pubnubPublishKey == "" || pubnubSubscribeKey == "" || pubnubSecretKey == "" || port == "" {
		fmt.Println("Invalid config flags")
		os.Exit(666)
	}

	pn := messaging.NewPubnub(pubnubPublishKey, pubnubSubscribeKey, pubnubSecretKey, "", false, uuid.NewV1().String())
	im := NewImageManager(yahooClientId, yahooClientSecret, redisUrl)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send(&im, pn))
	http.ListenAndServe(":"+port, nil)
}
