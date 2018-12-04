package main

import (
	"fmt"
	pb "social_media_app-golang/package1"
	"time"

	"golang.org/x/net/context"
)

func userExists(uname string) bool {
	if isServerAlive() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := rpcCaller.UserExists(ctx, &pb.UserExistsRequest{Username: uname})
		if err == nil {
			return reply.Status
		}
		fmt.Println("Debug: userExists rpc returned false", err)
		return false
	} else {
		debugPrint("Debug: Primary server down, cant process requests")
		return false
	}
}
