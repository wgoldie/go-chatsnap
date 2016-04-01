package main

import (
	"fmt"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/pubnub/go/messaging"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/satori/go.uuid"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/gopkg.in/dietsche/textbelt.v1"
	"net/http"
	"os"
	"strings"
)

func main() {
	bingAppId := os.Getenv("BING_APP_ID")
	pubnubPublishKey := os.Getenv("PUBNUB_PUBLISH_KEY")
	pubnubSubscribeKey := os.Getenv("PUBNUB_SUBSCRIBE_KEY")
	pubnubSecretKey := os.Getenv("PUBNUB_SECRET_KEY")
	redisUrl := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")
	smsTarget := os.Getenv("SMS_TARGET")

	if bingAppId == "" || pubnubPublishKey == "" || pubnubSubscribeKey == "" || pubnubSecretKey == "" || port == "" {
		fmt.Println("Invalid config flags!!!")
		os.Exit(666)
	}
	fmt.Println("a")
	pn := messaging.NewPubnub(pubnubPublishKey, pubnubSubscribeKey, pubnubSecretKey, "", false, uuid.NewV1().String())
	bm := NewBingManager(bingAppId, redisUrl)
	tb := TextBeltManager{
		Client:  textbelt.NewClientFromURL(textbelt.TextbeltAPIcanada),
		Numbers: strings.Split(smsTarget, " "),
	}
	fmt.Println("b")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/send", send(&bm, pn, tb))
	http.ListenAndServe(":"+port, nil)
	fmt.Println("c")
}
