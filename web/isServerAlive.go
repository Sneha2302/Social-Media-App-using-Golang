package main

import (
	"fmt"
	pb "social_media_app-golang/package1"
	"time"

	"golang.org/x/net/context"
)

func isServerAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := rpcCaller.HeartBeat(ctx, &pb.HeartBeatRequest{})
	if err == nil && reply.IsAlive {
		debugPrint("Debug: Heartbeat to Primary Successful")

		//Re-writing the FE servers global 'currentView' variable to make sure it matches with the Backend server
		currentView = int(reply.CurrentView)
		return true
	} else {
		debugPrint("Debug: Heartbeat to Primary Failed")
		//TODO: Start View Change here
		viewchange := &pb.PromptViewChangeArgs{NewView: int32(currentView + 1)}
		newprimary := GetPrimary(currentView+1, len(peers))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := peerRPC[newprimary].PromptViewChange(ctx, viewchange)
		if err == nil && reply.Success == true {
			currentView = currentView + 1
			rpcCaller = peerRPC[newprimary]
			primaryServerIndex = newprimary
			fmt.Println("Debug: we have a new primary")
			return true
		}

		//TODO: the below return false can be changed to true (if view change is successful)
		return false
	}
	return false
}
