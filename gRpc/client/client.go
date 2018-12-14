package client

import (
	"context"
	"google.golang.org/grpc"
	pb "twitterPt2/gRpc/protobuff"
	"log"
	"twitterPt2/storage"
	"time"
)

const (
	address = "localhost:9091"
)

func GetUserRpc(uName string) storage.Twitter_User {
		log.Printf("we in GetUserRpc")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.GetUser(ctx, &pb.GetUserRequest{Uname: uName})
	if err != nil {
		log.Printf("failed to call: %v", err)
	}
	user := r.Userinfo
		log.Printf("cp2")
	tmp := storage.ProtoToUser(user)
	return tmp
}

func TweetRpc(uName string, usr storage.Twitter_User) bool {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	pbUser := storage.UserToProto(usr)
	r, err := c.Tweet(ctx, &pb.TweetRequest{UName: uName, ProtoUser: pbUser})
	return r.T
}

func RegisterRpc(uName string, pWord1 string, pWord2 string) bool {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.AddUser(ctx, &pb.AddUserRequest{UName: uName, Pwd1: pWord1, Pwd2: pWord2})
	log.Printf("-------> RegisterRpc", r.T)
	return r.T

}

func UserExistsRpc(uName string, pWord string) bool {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.UserExists(ctx, &pb.UserExistsRequest{UName: uName, Pwd: pWord})
	log.Printf("-------> UserExistsRpc", r.T)
	return r.T
}

func GetUserInfoRpc(uName string) storage.Page {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.GetUserInfo(ctx, &pb.GetUserInfoRequest{Username: uName})
	log.Println("--------> TwitterPage", r)
	var twit = storage.Page{UName: r.UName, Unfollow: r.Unfollow, Follow: r.Follow, Posts: r.Posts}
	return twit
}

func FollowUserRpc(uName string, person string) bool {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.Follow(ctx, &pb.FollowRequest{UName: uName, Othername: person})
	return r.T
}

func UnfollowUserRpc(uName string, person string) bool {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWebClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	r, err := c.UnFollow(ctx, &pb.FollowRequest{UName: uName, Othername: person})
	return r.T
}
