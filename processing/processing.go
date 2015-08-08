package processing

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Identify(imgurl string, apikey string) (string, float64) {
	v := url.Values{}
	v.Set("url", imgurl)
	v.Set("outputMode", "json")
	v.Set("apikey", apikey)

	response, err := http.Get(fmt.Sprintf("https://access.alchemyapi.com/calls/url/URLGetRankedImageKeywords?%s", v.Encode()))

	if err != nil {
		log.Println(err)
	} else {
		log.Println(response)
	}

	os.Exit(1)

	return "true", 0
}
