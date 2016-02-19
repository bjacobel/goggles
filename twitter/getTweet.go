package twitter

import (
	"github.com/ChimeraCoder/anaconda"
)

// HasMedia is true if the tweet object has media entities
func HasMedia(tweet anaconda.Tweet) bool {
	return tweet.Entities.Media != nil
}

// GetMediaURL gets the URL of the first media entity in a tweet
func GetMediaURL(tweet anaconda.Tweet) string {
	return tweet.Entities.Media[0].Media_url_https
}

// Exclude filters certain things that are common and boring,
// e.g. screenshots get classified as "web site"
func Exclude(noun string) bool {

	switch noun {
	case "web site":
		return true
	case "person":
		return true
	default:
		return false
	}
}
