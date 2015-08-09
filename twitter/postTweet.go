package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"math/rand"
	"net/url"
	"time"
)

type Quip struct {
	Message    string
	Articleize bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
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
	Quip{"Does this look like %s to you? I think it is, but I'm never sure.", true},
	Quip{"Awesome, %s! Thanks for sharing ;)", true},
	Quip{"Coolest %s I've seen today!", false},
	Quip{"Love this %s!", false},
	Quip{"Wow, look at this %s! Super cool.", false},
	Quip{"This looks just like the %s I have!", false},
	Quip{"This might be my favorite %s I've ever seen.", false},
	Quip{"I want %s like this :(((", true},
	Quip{"Wonder where they got this sweet %s?", false},
	Quip{"Whoa, %s? Now I've seen everything.", true},
}

func Respond(tweet anaconda.Tweet, classification string, twitter anaconda.TwitterApi) bool {
	msg := fmt.Sprintf(
		"%s https://twitter.com/%s/status/%s",
		message(classification),
		tweet.User.IdStr,
		tweet.IdStr,
	)

	reply, errRpl := twitter.PostTweet(msg, url.Values{})

	if errRpl != nil {
		log.Printf("Error responding: ", errRpl)
		return false
	} else {
		log.Printf("Tweeted; link: https://twitter.com/%s/status/%d\n", reply.User.ScreenName, reply.Id)
		return true
	}
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
