package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/apex/go-apex"
	"gopkg.in/yaml.v2"

	"./alchemy"
	"./twitter"
)

// A Config struct holds configuration secrets
type Config struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
	AlchemyAPIKey     string `yaml:"alchemy_api_key"`
}

type message struct {
	Status string `json:"status"`
}

var twapi *anaconda.TwitterApi
var config Config

func init() {
	config = Parse("secrets.yml")

	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)

	twapi = anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		v := url.Values{}
		stream := twapi.PublicStreamSample(v)
		for {
			select {
			case tweet := <-stream.C:
				// Pull a tweet out of the channel, process it
				if handleTweet(tweet) == true {
					return &message{"success"}, nil
				}
			}
		}
	})
}

func handleTweet(t interface{}) bool {
	// Type assertion to anaconda.Tweet from interface{}
	if tweet, ok := t.(anaconda.Tweet); ok {
		if twitter.HasMedia(tweet) {
			classification, confidence := alchemy.Identify(twitter.GetMediaURL(tweet), config.AlchemyAPIKey)

			if confidence > 0.5 && !twitter.Exclude(classification) {
				return twitter.Respond(tweet, classification, *twapi)
			}
		}
	}

	return false
}

// Parse reads the config yaml file
func Parse(path string) Config {
	data, err := ioutil.ReadFile(path)

	if err == nil {
		var config Config

		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Fatal(err)
		}

		return config
	}

	// Else, could not read ./secrets.yml. Try to get from env vars
	config := Config{
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		AlchemyAPIKey:     os.Getenv("ALCHEMY_API_KEY"),
	}

	if config.ConsumerKey == "" || config.ConsumerSecret == "" || config.AccessToken == "" || config.AccessTokenSecret == "" || config.AlchemyAPIKey == "" {
		log.Fatal("Missing secrets.yaml *and* environment variables")
		return Config{}
	}

	// Else
	return config
}
