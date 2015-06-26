package twitter

import (
	"github.com/chimeracoder/anaconda"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
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

func Stream() {
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
			go HandleTweet(tweet)
		}
	}
}

func HandleTweet(tweet interface{}) {
	log.Println("%s", spw.Sdump(tweet))
}
