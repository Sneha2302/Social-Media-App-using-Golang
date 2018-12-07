package main

import (
	"fmt"
	"time"

	//pb "social_media_app-golang/package1"
	pb "package1"

	"golang.org/x/net/context"
)

func deleteUser(username string) int {
	//TODO: for later stages, we might have to add Locks here
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := rpcCaller.DeleteUser(ctx, &pb.Credentials{Uname: username})
	if err == nil {
		fmt.Println("Delete User RPC successful", reply)
		return 0
	} else {
		fmt.Println("Delete User RPC failed", reply, err)
		return -1
	}
}
