package main

import (
	"fmt"
	//pb "social_media_app-golang/package1"
	pb "package1"
	"time"

	"golang.org/x/net/context"
)

func userExists(uname string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := rpcCaller.UserExists(ctx, &pb.UserExistsRequest{Username: uname})
	if err == nil {
		return reply.Status
	}
	fmt.Println("Debug: userExists rpc returned false", err)
	return false
}
