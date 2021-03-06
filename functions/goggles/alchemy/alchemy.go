package alchemy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type imageKeyword struct {
	Text  string `json: "text"`
	Score string `json: "score"`
}

type response struct {
	Status            string         `json: "status"`
	URL               string         `json: "url"`
	TotalTransactions string         `json: "totalTransactions"`
	ImageKeywords     []imageKeyword `json: "imageKeywords"`
}

// Identify runs an image url against the Alchemy API
func Identify(imgurl string, apikey string) (string, float64) {
	v := url.Values{}
	v.Set("url", imgurl)
	v.Set("outputMode", "json")
	v.Set("apikey", apikey)

	responseJSON, httpErr := http.Get(fmt.Sprintf("https://access.alchemyapi.com/calls/url/URLGetRankedImageKeywords?%s", v.Encode()))
	defer responseJSON.Body.Close()

	if httpErr != nil {
		log.Println(httpErr)
	} else {
		body, ioErr := ioutil.ReadAll(responseJSON.Body)
		if ioErr != nil {
			log.Println(ioErr)
		} else {
			var response response

			if marshalErr := json.Unmarshal(body, &response); marshalErr != nil {
				log.Println(marshalErr)
			} else {
				if len(response.ImageKeywords) > 0 {
					score, _ := strconv.ParseFloat(response.ImageKeywords[0].Score, 64)
					return response.ImageKeywords[0].Text, score
				}
			}
		}
	}

	return "failure", 0
}
