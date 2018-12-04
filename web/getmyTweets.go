package main

import (
	"fmt"
	pb "social_media_app-golang/package1"
	"time"

	"golang.org/x/net/context"
)

func getMyTweets(username string) *pb.OwnTweetsReply {
	if isServerAlive() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := rpcCaller.OwnTweets(ctx, &pb.OwnTweetsRequest{Username: username})
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return reply
	} else {
		debugPrint("Debug: Primary server down, cant process requests")
		return nil
	}
}
