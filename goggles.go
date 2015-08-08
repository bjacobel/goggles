package main

import (
	"github.com/bjacobel/goggles/twitter"
	"github.com/chimeracoder/anaconda"
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
}

func main() {
	v := url.Values{}
	v.Set("language", "en")

	config := Parse("secrets.yml")
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)

	twapi := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	stream := twapi.PublicStreamSample(v)

	for {
		select {
		case <-stream.Quit:
			log.Fatal("Stream terminated, wrapping up and quitting...")
			break
		case tweet := <-stream.C:
			// Pull a tweet out of the channel, process it (currently this is synchronous)
			twitter.HandleTweet(tweet, *twapi)
		}
	}
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
