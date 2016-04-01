package main

import (
	"fmt"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/pubnub/go/messaging"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/satori/go.uuid"
	"net/http"
	"os"
)

func main() {
	bingAppId := os.Getenv("BING_APP_ID")
	pubnubPublishKey := os.Getenv("PUBNUB_PUBLISH_KEY")
	pubnubSubscribeKey := os.Getenv("PUBNUB_SUBSCRIBE_KEY")
	pubnubSecretKey := os.Getenv("PUBNUB_SECRET_KEY")
	redisUrl := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")

	if bingAppId == "" || pubnubPublishKey == "" || pubnubSubscribeKey == "" || pubnubSecretKey == "" || port == "" {
		fmt.Println("Invalid config flags!!!")
		os.Exit(666)
	}
	fmt.Println("a")
	pn := messaging.NewPubnub(pubnubPublishKey, pubnubSubscribeKey, pubnubSecretKey, "", false, uuid.NewV1().String())
	bm := NewBingManager(bingAppId, redisUrl)
	fmt.Println("b")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send(&bm, pn))
	http.ListenAndServe(":"+port, nil)
	fmt.Println("c")
}
