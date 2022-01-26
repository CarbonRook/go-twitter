package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/carbonrook/go-twitter/twitter"
)

func main() {

	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	if bearerToken == "" {
		fmt.Printf("TWITTER_BEARER_TOKEN environment variable not set")
		os.Exit(0)
	}

	// Twitter Client
	client := twitter.NewClientWithBearer(http.DefaultClient, bearerToken)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}
	demux.StreamData = func(sd *twitter.StreamData) {
		dataJson, err := json.MarshalIndent(sd, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal StreamData: %s", err)
		} else {
			fmt.Printf("StreamData: %s\n", dataJson)
		}
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	filterParams := &twitter.StreamParams{
		TweetFields: []string{"created_at", "text", "id", "lang", "public_metrics", "referenced_tweets", "conversation_id", "entities"},
		Expansions:  []string{"author_id", "referenced_tweets.id", "in_reply_to_user_id", "entities.mentions.username", "referenced_tweets.id.author_id"},
		UserFields:  []string{"id", "name", "username", "public_metrics", "verified"},
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
