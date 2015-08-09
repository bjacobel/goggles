package twitter

import (
	"github.com/ChimeraCoder/anaconda"
)

func HasMedia(tweet anaconda.Tweet) bool {
	return tweet.Entities.Media != nil
}

func GetMediaURL(tweet anaconda.Tweet) string {
	return tweet.Entities.Media[0].Media_url_https
}

func Exclude(noun string) bool {
	// Certain things are common and boring, e.g. screenshots get classified as "web site"
	// Exclude (don't tweet) these ones.
	switch noun {
	case "web site":
		return true
	case "person":
		return true
	default:
		return false
	}
}
