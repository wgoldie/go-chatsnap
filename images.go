package main

import (
	"github.com/mrjones/oauth"
    "encoding/json"
)

type ImageManager struct {
	Url		string
	Consumer 	*oauth.Consumer
}

func (im *ImageManager) getImageUrl(query string) string {
	queryString := map[string]string{
		"q": "\"" + query + "\"",
		"filter":"no",
		"dimensions": "small",
		"count":"1"}

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
	   
	return m.BossResponse.Images.Results[0].Url
}

func NewImageManager(ClientId string, ClientSecret string) ImageManager {
	return ImageManager{
		Url: "https://yboss.yahooapis.com/ysearch/web",
		Consumer: oauth.NewConsumer(ClientId, ClientSecret, oauth.ServiceProvider{})}
}