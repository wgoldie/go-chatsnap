package main

import (
	"encoding/json"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/mrjones/oauth"
	"github.com/go-redis/redis"
)

type ImageManager struct {
	Url      string
	Consumer *oauth.Consumer
	Client *redis.Client
}

func (im *ImageManager) getImageUrl(query string) string {
	queryString := map[string]string{
		"q":          "\"" + query + "\"",
		"filter":     "no",
		"dimensions": "small",
		"count":      "1"}

	accessToken := &oauth.AccessToken{}
	r, err := im.Consumer.Get("https://yboss.yahooapis.com/ysearch/images", queryString, accessToken)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var m struct {
		BossResponse struct {
			Images struct {
				Results []struct {
					Url string `json:"url"`
				} `json:"results"`
			} `json:"images"`
		} `json:"bossresponse"`
	}

	err = decoder.Decode(&m)
	if err != nil {
		panic(err)
	}

	url := m.BossResponse.Images.Results[0].Url
	err = im.Client.Set(query, url)
	return url
}

func NewImageManager(ClientId string, ClientSecret string, RedisUrl string) ImageManager {
	return ImageManager{
		Url:      "https://yboss.yahooapis.com/ysearch/web",
		Consumer: oauth.NewConsumer(ClientId, ClientSecret, oauth.ServiceProvider{})
		Client: redis.NewClient(&redis.Options{
			Addr: redisUrl,
			Password: "",
			DB: 0,
		})}
}
