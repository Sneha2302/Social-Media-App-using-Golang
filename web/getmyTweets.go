package main

import (
	"fmt"
	//pb "social_media_app-golang/package1"
	pb "package1"
	"time"

	"golang.org/x/net/context"
)

func getMyTweets(username string) *pb.OwnTweetsReply {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := rpcCaller.OwnTweets(ctx, &pb.OwnTweetsRequest{Username: username})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return reply
}
