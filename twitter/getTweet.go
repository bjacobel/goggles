package twitter

import (
	"github.com/bjacobel/goggles/processing"
	"github.com/chimeracoder/anaconda"
	"log"
)

func HandleTweet(t interface{}, twitter anaconda.TwitterApi) {
	// Type assertion to anaconda.Tweet from interface{}
	if tweet, ok := t.(anaconda.Tweet); ok {
		if hasMedia(tweet) {
			classification, confidence := processing.Identify(getMediaURL(tweet))

			if confidence > 0.5 && !exclude(classification[0]) {
				respond(tweet, classification, twitter)
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

func exclude(noun string) bool {
	// Certain things are common and boring, e.g. screenshots get classified as "web site"
	// Exclude (don't tweet) these ones.
	switch noun {
	case "web site":
		return true
	default:
		return false
	}
}
