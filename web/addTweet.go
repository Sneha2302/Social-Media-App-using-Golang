package main

import (
	"fmt"
	"time"

	pb "social_media_app-golang/package1"

	"golang.org/x/net/context"
)

func addTweet(username string, tweettext string) {
	if isServerAlive() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := rpcCaller.AddTweet(ctx, &pb.AddTweetRequest{Username: username, TweetText: tweettext, Broadcast: true})
		if err != nil {
			fmt.Println("Debug: tweet addition failed", err)
		}
	} else {
		debugPrint("Debug: Primary server down, cant process requests")
	}
}
