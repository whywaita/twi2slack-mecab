package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/nlopes/slack"
	mecab "github.com/shogo82148/go-mecab"
)

func validationEnviroments() error {
	if os.Getenv("TWITTER_CONSUMER_KEY") == "" {
		return errors.New("TWITTER_CONSUMER_KEY is blank")
	}
	if os.Getenv("TWITTER_CONSUMER_SECRET") == "" {
		return errors.New("TWITTER_CONSUMER_KEY is blank")
	}
	if os.Getenv("TWITTER_OAUTH_TOKEN") == "" {
		return errors.New("TWITTER_OAUTH_TOKEN is blank")
	}
	if os.Getenv("TWITTER_OAUTH_TOKEN_SECRET") == "" {
		return errors.New("TWITTER_OAUTH_TOKEN_SECRET is blank")
	}
	if os.Getenv("SLACK_TOKEN") == "" {
		return errors.New("SLACK_TOKEN is blank")
	}
	if os.Getenv("SLACK_CHANNEL") == "" {
		return errors.New("SLACK_CHANNEL is blank")
	}

	return nil
}

func postSlack(api *slack.Client, channelName string, message string) error {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		//Pretext: "pretext",
		Text: message,
	}

	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(channelName, "Twitter Search Result", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func initialize() (*anaconda.TwitterApi, *slack.Client, string, error) {
	// validation Enviroments value
	err := validationEnviroments()
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, nil, "", err
	}

	// set Twitter infomation
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	twitterAPI := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"))
	twitterAPI.SetLogger(anaconda.BasicLogger)

	slackAPI := slack.New(os.Getenv("SLACK_TOKEN"))
	slackChannel := os.Getenv("SLACK_CHANNEL")

	return twitterAPI, slackAPI, slackChannel, nil
}

func main() {
	var err error

	twitterAPI, slackAPI, slackChannel, err := initialize()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// make mecab map
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	// Twitter connect
	v := url.Values{}
	stream := twitterAPI.UserStream(v)

	// User Stream wait...
	for {
		select {
		case item := <-stream.C:
			switch status := item.(type) {
			case anaconda.Tweet:
				if !strings.Contains(status.Text, "わいわいた") {
				} else {
					// get Tweet
					result, err := tagger.Parse(status.Text)
					if err != nil {
						panic(err)
					}
					postSlack(slackAPI, slackChannel, result)
				}
			default:
				// nothing
			}
		}
	}

}
