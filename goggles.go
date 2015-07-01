package main

import (
	//"fmt"
	"github.com/chimeracoder/anaconda"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var spw = spew.NewDefaultConfig()

func init() {
	// Add some configuration to our JSON logger
	spw = &spew.ConfigState{Indent: "\t", MaxDepth: 5}
}

type Config struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
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

func main() {
	config := Parse("secrets.yml")

	// can't stop giggling every time I see "anaconda"
	// thank you, sir mix-a-lot
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	twitter := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)

	v := url.Values{}
	v.Set("language", "en")

	stream := twitter.PublicStreamSample(v)

	for {
		select {
		case <-stream.Quit:
			log.Fatal("Stream terminated, wrapping up and quitting...")
			break
		case tweet := <-stream.C:
			go handleTweet(tweet)
		}
	}
}

func handleTweet(t interface{}) {
	// Type assertion to anaconda.Tweet from interface{}
	if tweet, ok := t.(anaconda.Tweet); ok {
		if hasMedia(tweet) {
			identify(getMediaURL(tweet))
		}
	} else {
		log.Fatal("Tried to handle a tweet that was not a tweet (?)")
	}
}

func hasMedia(tweet anaconda.Tweet) bool {
	return tweet.Entities.Media != nil
}

func getMediaURL(tweet anaconda.Tweet) string {
	return tweet.Entities.Media[0].Media_url_https
}

func identify(url string) {
	// download it, save to tmp
	fileName := download(url)

	// wd, _ := os.Getwd()

	// process it
	out, err := exec.Command("overfeat/bin/linux_32/overfeat", "-n 1", fileName).CombinedOutput()

	if err != nil {
		log.Fatal(out, err)
	} else {
		log.Printf("%s\n", out)
	}

	// delete from /tmp
}

func download(url string) string {
	tokens := strings.Split(url, "/")
	fileName := "/tmp/goggles/" + tokens[len(tokens)-1]

	output, err := os.Create(fileName)

	if err != nil {
		log.Fatal("Error while creating", fileName, "-", err)
	}

	defer output.Close()

	response, err := http.Get(url)

	if err != nil {
		log.Fatal("Error while downloading", url, "-", err)
	}

	defer response.Body.Close()

	_, er := io.Copy(output, response.Body)

	if er != nil {
		log.Fatal("Error while saving", url, "-", er)
	}

	return fileName
}
