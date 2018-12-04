package main

import (
	"fmt"
	"time"

	pb "social_media_app-golang/package1"

	"golang.org/x/net/context"
)

func deleteUser(username string) int {
	if isServerAlive() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := rpcCaller.DeleteUser(ctx, &pb.Credentials{Uname: username, Broadcast: true})
		if err == nil {
			fmt.Println("Delete User RPC successful", reply)
			return 0
		} else {
			fmt.Println("Delete User RPC failed", reply, err)
			return -1
		}
	} else {
		debugPrint("Debug: Primary server down, cant process requests")
		return -1
	}
}
