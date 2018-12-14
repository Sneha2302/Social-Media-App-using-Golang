package main

import (
	"context"
	"log"
	"net"
	"sync"
	pb "twitterPt2/gRpc/protobuff"
	"twitterPt2/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9091"
)

type DB struct {
	mu        sync.Mutex
	UsersInfo map[string]storage.Twitter_User
}

func (db *DB) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	var uName string = in.Uname
	var temp storage.Twitter_User = db.UsersInfo[uName]
	user := storage.UserToProto(temp)
	log.Printf("------> server user", user)
	return &pb.GetUserReply{Userinfo: user}, nil
}

func (db *DB) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.BoolReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	uName := in.UName   
	password1 := in.Pwd1 
	password2 := in.Pwd2 

	if password1 != password2 {
		return &pb.BoolReply{T: false}, nil
	}
	if uName == "" || password1 == "" {
		return &pb.BoolReply{T: false}, nil
	}	
	log.Printf("......adding user........")
	StoreUser := storage.Twitter_User{uName, password1, storage.Tweets{}, []string{uName}}
	if _, ok := db.UsersInfo[uName]; ok {
		return &pb.BoolReply{T: false}, nil
	}
	db.UsersInfo[uName] = StoreUser
	return &pb.BoolReply{T: true}, nil
}

func (db *DB) Tweet(ctx context.Context, in *pb.TweetRequest) (*pb.BoolReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	uName := in.UName
	usr := storage.ProtoToUser(in.ProtoUser)
	if uName != usr.UName {
		return &pb.BoolReply{T: false}, nil
	}
	if _, ok := db.UsersInfo[uName]; ok != true {
		return &pb.BoolReply{T: false}, nil
	}
	db.UsersInfo[uName] = usr
	return &pb.BoolReply{T: true}, nil
}

func (db *DB) UserExists(ctx context.Context, in *pb.UserExistsRequest) (*pb.BoolReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	uName := in.UName    
	password := in.Pwd 
	if uName == "" || password == "" {
		return &pb.BoolReply{T: false}, nil
	}
	user, exist := db.UsersInfo[uName]
	log.Printf("-------> user ", user)
	log.Printf("-------> exists ", exist)
	if exist && user.Pwd == password {
		return &pb.BoolReply{T: true}, nil
	}
	return &pb.BoolReply{T: false}, nil
}

func (db *DB) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.BoolReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	uName := in.UName     
	otherName := in.Othername 
	if twitter_user, ok := db.UsersInfo[uName]; ok {
		if storage.Contains(twitter_user.Follow, otherName) {
			return &pb.BoolReply{T: false}, nil
		}
		twitter_user.Follow = append(twitter_user.Follow, otherName)
		db.UsersInfo[uName] = twitter_user
		return &pb.BoolReply{T: true}, nil
	}
	return &pb.BoolReply{T: false}, nil
}

func (db *DB) UnFollow(ctx context.Context, in *pb.FollowRequest) (*pb.BoolReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	uName := in.UName      //form
	otherName := in.Othername //form
	if twitter_user, ok := db.UsersInfo[uName]; ok {
		if !storage.Contains(twitter_user.Follow, otherName) {
			return &pb.BoolReply{T: false}, nil
		}
		twitter_user.Follow = storage.Deletes(twitter_user.Follow, otherName)
		db.UsersInfo[uName] = twitter_user
		return &pb.BoolReply{T: true}, nil
	}
	return &pb.BoolReply{T: false}, nil
}


func GetTweets(arr storage.Tweets) []string {
	var ret []string
	for _, tweet := range arr {
		tmp := tweet.User + ": " + tweet.Contents
		ret = append(ret, tmp)
	}
	return ret
}

func (db *DB) GetUserInfo(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	uName := in.Username //form
	user, _ := db.UsersInfo[uName]
	log.Printf("-------> TwitterPage Userinfo ", user)
	UserName := user.UName
	Following := user.Follow
	log.Printf("..............", Following)
	var UnFollowed []string
	var Posts storage.Tweets
	for name, userInfo := range db.UsersInfo {
		if storage.Contains(Following, name) {
			for _, post := range userInfo.Post {
				Posts = append(Posts, post)
			}
		} else {
			UnFollowed = append(UnFollowed, name)
		}
	}
	userPosts := storage.GetUserInfo(Posts)
	Following = storage.Deletes(Following, uName)
	var tweet = &pb.Page{UName: UserName, Unfollow: UnFollowed, Follow: Following, Posts: userPosts}
	return &pb.GetUserInfoReply{UName: tweet.UName, Unfollow: tweet.Unfollow, Follow: tweet.Follow, Posts: tweet.Posts}, nil

}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	WebDB := &DB{}
	WebDB.UsersInfo = make(map[string]storage.Twitter_User)
	pb.RegisterWebServer(s, WebDB)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
