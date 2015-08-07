package main

import (
	"fmt"
	"github.com/chimeracoder/anaconda"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var spw = spew.NewDefaultConfig()
var twitter *anaconda.TwitterApi

func init() {
	// Add some configuration to our JSON logger
	spw = &spew.ConfigState{Indent: "\t", MaxDepth: 5}

	// Parse a YAML file to get secret configuration
	config := Parse("secrets.yml")

	// can't stop giggling every time I see "anaconda"
	// thank you, sir mix-a-lot
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	twitter = anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
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
	v := url.Values{}
	v.Set("language", "en")

	stream := twitter.PublicStreamSample(v)

	for {
		select {
		case <-stream.Quit:
			log.Fatal("Stream terminated, wrapping up and quitting...")
			break
		case tweet := <-stream.C:
			handleTweet(tweet)
		}
	}
}

func handleTweet(t interface{}) {
	// Type assertion to anaconda.Tweet from interface{}
	if tweet, ok := t.(anaconda.Tweet); ok {
		if hasMedia(tweet) {
			classification, confidence := identify(getMediaURL(tweet))

			if confidence > 0.5 {
				respond(tweet, classification)
			}
		}
	} else {
		log.Println("Tried to handle a tweet that was not a tweet (?)")
	}
}

func hasMedia(tweet anaconda.Tweet) bool {
	return tweet.Entities.Media != nil
}

func getMediaURL(tweet anaconda.Tweet) string {
	return tweet.Entities.Media[0].Media_url_https
}

func identify(url string) ([]string, float64) {
	// download it, save to tmp
	fileName := download(url)

	// process it
	out, err := exec.Command("./OverFeat/src/overfeat", "-n 1", fileName).CombinedOutput()

	if err != nil {
		log.Println(out, err)
	}

	split := strings.Split(string(out), " ")
	confidence, _ := strconv.ParseFloat(strings.Trim(split[len(split)-1], "\n"), 64)
	classifier := strings.Split(strings.Trim(strings.Join(split[:len(split)-1], " "), ","), ", ")

	// delete from /tmp

	return classifier, confidence
}

func download(url string) string {
	tokens := strings.Split(url, "/")
	fileName := "/tmp/" + tokens[len(tokens)-1]

	output, err := os.Create(fileName)

	if err != nil {
		log.Println("Error while creating", fileName, " - ", err)
	}

	defer output.Close()

	response, err := http.Get(url)

	if err != nil {
		log.Println("Error while downloading", url, " - ", err)
	}

	defer response.Body.Close()

	_, er := io.Copy(output, response.Body)

	if er != nil {
		log.Println("Error while saving", url, "-", er)
	}

	return fileName
}

func respond(tweet anaconda.Tweet, classification []string) {
	msg := fmt.Sprintf(
		"%s https://twitter.com/%s/status/%s",
		message(classification[0]),
		tweet.User.IdStr,
		tweet.IdStr,
	)

	reply, errRpl := twitter.PostTweet(msg, url.Values{})

	if errRpl != nil {
		log.Printf("Error responding: ", errRpl)
		return
	} else {
		log.Printf("Tweeted; link: https://twitter.com/%s/status/%d\n", reply.User.ScreenName, reply.Id)
	}
}

type Quip struct {
	Message    string
	Articleize bool
}

// pairs of messages and whether the noun must be articleized to fit them
var quipOptions = []Quip{
	Quip{"Look what I found! I think it's %s.", true},
	Quip{"This looks like %s!", true},
	Quip{"What do we have here? It looks like %s.", true},
	Quip{"I found %s!", true},
	Quip{"I've never seen %s like this before!", true},
	Quip{"Wow, I almost didn't recognize this as %s!", true},
	Quip{"Cool! It's %s!", true},
	Quip{"I could be wrong, but this looks like %s to me.", true},
	Quip{"Bet you didn't think I'd recognize this ;) It's %s!", true},
	Quip{"Wow, %s! Didn't think I'd see one of those today :)", true},
	Quip{"I found %s! Now I've seen everything :)", true},
	Quip{"Look, %s!", true},
	Quip{"Does this look like %s to you? I think it is, but I'm never sure!", true},
	Quip{"Awesome, %s! Thanks for sharing ;)", true},
	Quip{"Coolest %s I've seen today!", false},
	Quip{"Love this %s!", false},
	Quip{"Wow, look at this %s! Super cool.", false},
	Quip{"This looks just like the %s I have!", false},
	Quip{"This might be my favorite %s I've ever seen.", false}
	Quip{"I want %s like this :(((", true},
	Quip{"Wonder where they got this sweet %s?", false},
	Quip{"Whoa, %s? Now I've seen everything.", true},
}

func message(noun string) string {
	choice := quipOptions[rand.Intn(len(quipOptions))]

	if choice.Articleize {
		return fmt.Sprintf(choice.Message, articleize(noun))
	} else {
		return fmt.Sprintf(choice.Message, noun)
	}
}

func isVowel(letter string) bool {
	// @TODO: This is super quick n' dirty. Obviously English has lots of corner cases
	// not covered by this, we want to test for vowel "sounds" not just vowels
	switch letter {
	case "a":
		return true
	case "e":
		return true
	case "i":
		return true
	case "o":
		return true
	case "u":
		return true
	default:
		return false
	}
}

func articleize(noun string) string {
	if isVowel(noun[:1]) {
		return fmt.Sprintf("an %s", noun)
	} else {
		return fmt.Sprintf("a %s", noun)
	}
}
