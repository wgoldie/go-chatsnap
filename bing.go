package main

import (
    "encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"net/http"
	"net/url"
	"strings"
    "io/ioutil"
)

// Holds information on the current image API
type BingManager struct {
	AccountKey string
	Client     *redis.Client
	HTTPClient *http.Client
}

// Queries the image search api for a new image url for the given query string
func (bm *BingManager) queryNewImageUrl(query string) (string, error) {
	queryString := fmt.Sprintf(
        "https://api.datamarket.azure.com/Bing/Search/Image?$format=json&Query=%%27{%s}%%27&$top=1",
        query)

	req, err := http.NewRequest("GET", queryString, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(bm.AccountKey, bm.AccountKey)

	resp, err := bm.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return "", err
    }
    
    var m struct{
        D struct { 
            Results []struct {
                MediaUrl string `json:"MediaUrl"`
            } `json:"results"`
        } `json:"d"`
    }

    err = json.Unmarshal(body, &m)
	if err != nil {
        fmt.Println(err)
		return "", err
	}
    return m.D.Results[0].MediaUrl, err
}

// Queries the cached image url database for the given query string
func (bm *BingManager) queryCachedImageUrl(query string) (bool, string, error) {
	val, err := bm.Client.Get(query).Result()

	if err == redis.Nil {
		return false, "", nil
	} else if err != nil {
		return false, val, nil
	}

	//	fmt.Println("Cache returned result")

	return true, val, nil
}

// Retrieves the url for the given query string if it is present in the cached image url databse
// Otherwise, queries for a new url and caches it in the database
func (bm *BingManager) getImageUrl(query string) (string, error) {
	found, val, err := bm.queryCachedImageUrl(query)

	if err != nil {
		return "", err
	}

	if !found {
		val, err = bm.queryNewImageUrl(query)
		if err != nil || val == "" {
			return "", err
		}

		// todo: check for race condition replacement?
		// no real reason to do so...
		// it's unlikely that an image would be double set fast enough
		// and would also have a different result from the search api
		// and it wouldn't matter
		err = bm.Client.Set(query, val, 0).Err()
	}

	return val, err
}

// Gets a series of imageurls for the given query string's ngram elements
// Currently seperates on spaces
func (bm *BingManager) getImageUrls(query string) []string {
    fmt.Println(query)
	fields := strings.Fields(query)
	results := []string{}
	for _, el := range fields {
		val, err := bm.getImageUrl(el)

		if err != nil {
			fmt.Println(err)
			continue
		} else if val != "" {
			results = append(results, val)
		}
	}

	return results
}

// Constructs a new BingManager struct
func NewBingManager(accountKey string, RedisUrl string) BingManager {

	parsedURL, _ := url.Parse(RedisUrl)
	password, _ := parsedURL.User.Password()
	host := parsedURL.Host

	return BingManager{
		AccountKey: accountKey,
		Client: redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       0,
		}),
		HTTPClient: &http.Client{},
	}
}
