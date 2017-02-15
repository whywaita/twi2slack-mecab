package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	mecab "github.com/shogo82148/go-mecab"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// set Twitter infomation
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"))
	api.SetLogger(anaconda.BasicLogger) // set logger

	// make mecab map
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	// Twitter connect
	v := url.Values{}
	stream := api.UserStream(v)

	// User Stream wait...
	for {
		select {
		case item := <-stream.C:
			switch status := item.(type) {
			case anaconda.Tweet:
				// get Tweet
				result, err := tagger.Parse(status.Text)
				if err != nil {
					panic(err)
				}

				fmt.Println(result)
			default:
			}
		}
	}

}
