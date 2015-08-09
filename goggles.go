package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/bjacobel/goggles/alchemy"
	"github.com/bjacobel/goggles/twitter"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
)

type Config struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
	AlchemyAPIKey     string `yaml:"alchemy_api_key"`
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
	v := url.Values{}
	v.Set("language", "en")

	c := cron.New()
	c.Start()
	defer c.Stop()

	c.AddFunc("@every 30m", func() {
		stream := twapi.PublicStreamSample(v)
		for {
			select {
			case <-stream.Quit:
				log.Fatal("Stream terminated, wrapping up and quitting...")
				break
			case tweet := <-stream.C:
				// Pull a tweet out of the channel, process it
				if handleTweet(tweet) == true {
					break
				}
			}
		}
	})

	select {}
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
	} else {
		log.Println("Tried to handle a tweet that was not a tweet (?)")
	}

	return false
}

func Parse(path string) Config {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
