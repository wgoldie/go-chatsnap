package main

import (
	"encoding/json"
	"fmt"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/github.com/mrjones/oauth"
	"github.com/wgoldie/go-chatsnap/Godeps/_workspace/src/gopkg.in/redis.v3"
	"net/url"
	"strings"
)

type ImageManager struct {
	Url      string
	Consumer *oauth.Consumer
	Client   *redis.Client
}

func (im *ImageManager) queryNewImageUrl(query string) (string, error) {
	queryString := map[string]string{
		"q":          "\"" + query + "\"",
		"filter":     "no",
		"dimensions": "small",
		"count":      "1"}

	accessToken := &oauth.AccessToken{}
	r, err := im.Consumer.Get(im.Url, queryString, accessToken)
	if err != nil {
		return "", err
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

	fmt.Println("Queried for new result")
	return m.BossResponse.Images.Results[0].Url, err
}

func (im *ImageManager) queryCachedImageUrl(query string) (bool, string, error) {
	val, err := im.Client.Get(query).Result()

	if err == redis.Nil {
		return false, "", nil
	} else if err != nil {
		return false, val, nil
	}

	fmt.Println("Cache returned result")

	return true, val, nil
}

func (im *ImageManager) getImageUrl(query string) (string, error) {
	found, val, err := im.queryCachedImageUrl(query)

	if err != nil {
		return "", err
	}

	if !found {
		val, err = im.queryNewImageUrl(query)
		if err != nil {
			return "", err
		}

		// todo: check for race condition replacement?
		// no real reason to do so...
		// it's unlikely that an image would be double set fast enough
		// and would also have a different result from the search api
		// and it wouldn't matter
		err = im.Client.Set(query, val, 0).Err()
	}

	return val, err
}

func (im *ImageManager) getImageUrls(query string) []string {
	fields := strings.Fields(query)
	results := []string{}
	for _, el := range fields {
		val, err := im.getImageUrl(el)

		if err != nil {
			fmt.Println(err)
			continue
		} else if val != "" {
			results = append(results, val)
		}
	}

	return results
}

func NewImageManager(ClientId string, ClientSecret string, RedisUrl string) ImageManager {

	parsedURL, _ := url.Parse(RedisUrl)
	password, _ := parsedURL.User.Password()
	host := parsedURL.Host

	return ImageManager{
		Url:      "https://yboss.yahooapis.com/ysearch/images",
		Consumer: oauth.NewConsumer(ClientId, ClientSecret, oauth.ServiceProvider{}),
		Client: redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       0,
		})}
}
