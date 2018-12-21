package client

import (
	"context"
	"fmt"
	"log"
	pb "social_media_app-golang/gRpc/protobuff"
	"social_media_app-golang/storage"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var curLeader int = -1
var idx int
var addresslist = []string{"localhost:9093", "localhost:9094", "localhost:9095"}

func RpcEstablish(addresslist []string) (*grpc.ClientConn, pb.WebClient) {
	if curLeader == -1 {
		idx = 0
	} else {
		idx = curLeader
	}
	for idx < len(addresslist) {
		address := addresslist[idx]
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c := pb.NewWebClient(conn)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		port := strings.Split(address, ":")[1]
		fmt.Println(port)
		var add = &pb.IsLeaderRequest{Address: port}
		t, err := c.IsLeader(ctx, add)
		log.Println("This is the leader node", t)
		log.Println("Error--------------------->", err)
		log.Println("------>leader", idx)
		if err == nil && t.T {
			log.Println("------> Leader", idx)
			curLeader = idx
			return conn, c
		} else {
			idx += 1
			if idx == len(addresslist) {
				idx = 0
			}
		}
		conn.Close()
	}
	log.Println("----> pointer")
	return nil, nil
}

func GetUserRpc(uName string) storage.Twitter_User {
	log.Printf("we in GetUserRpc")
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.GetUser(ctx, &pb.GetUserRequest{Uname: uName})
	if err != nil {
		log.Printf("Unable to call %v", err)
		// return nil
	}
	user := r.Userinfo
	log.Printf("We got user")
	tmp := storage.ProtoToUser(user)
	return tmp
}

func TweetRpc(uName string, usr storage.Twitter_User) bool {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewWebClient(conn)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	pbUser := storage.UserToProto(usr)
	r, _ := c.Tweet(ctx, &pb.TweetRequest{UName: uName, ProtoUser: pbUser})
	return r.T
}

func RegisterRpc(uName string, pWord1 string, pWord2 string) bool {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewWebClient(conn)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, _ := c.AddUser(ctx, &pb.AddUserRequest{UName: uName, Pwd1: pWord1, Pwd2: pWord2})
	//log.Printf("-------> RegisterRpc", r.T)
	return r.T

}

func UserExistsRpc(uName string, pWord string) bool {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewWebClient(conn)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, _ := c.UserExists(ctx, &pb.UserExistsRequest{UName: uName, Pwd: pWord})

	return r.T
}

func GetUserInfoRpc(uName string) storage.Page {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewWebClient(conn)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, _ := c.GetUserInfo(ctx, &pb.GetUserInfoRequest{Username: uName})
	log.Println("--------> loading Twitter Page ------->", r)
	var tweet = storage.Page{UName: r.UName, Unfollow: r.Unfollow, Follow: r.Follow, Posts: r.Posts}
	return tweet
}

func FollowUserRpc(uName string, person string) bool {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewWebClient(conn)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, _ := c.Follow(ctx, &pb.FollowRequest{UName: uName, Othername: person})
	return r.T
}

func UnfollowUserRpc(uName string, person string) bool {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewWebClient(conn)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, _ := c.UnFollow(ctx, &pb.FollowRequest{UName: uName, Othername: person})
	return r.T
}
func RpcJoin(nodeID string, remoteAddr string) {
	// addresslist = append(addresslist, "local:host"+remoteAddr)
	conn, c := RpcEstablish(addresslist)
	defer conn.Close()
	log.Println("Joining")
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	var join = &pb.JoinRequest{NodeID: nodeID, RemoteAddr: remoteAddr}
	c.Join(ctx, join)
}
