package storage

import (
	"fmt"
	pb "social_media_app-golang/gRpc/protobuff"
)

type Twitter_User struct {
	UName  string
	Pwd    string
	Post   Tweets
	Follow []string
}

type Tweets []TwitterPosts

type Page struct {
	UName    string
	Follow   []string
	Unfollow []string
	Posts    []string
}

type TwitterPosts struct {
	Contents string
	Date     int64
	User     string
}

func ProtoToUser(user *pb.Twitter_User) Twitter_User {
	var posts Tweets
	fmt.Println(user)
	for _, post := range user.Posts {
		var tempPost = TwitterPosts{Contents: post.Contents, Date: post.Date, User: post.User}
		posts = append(posts, tempPost)
	}
	var temp = Twitter_User{UName: user.UName, Pwd: user.Pwd, Post: posts, Follow: user.Follow}
	return temp
}

func UserToProto(temp Twitter_User) *pb.Twitter_User {
	var posts []*pb.TwitterPosts
	for _, post := range temp.Post {
		var tempPost = &pb.TwitterPosts{Contents: post.Contents, Date: post.Date, User: post.User}
		posts = append(posts, tempPost)
	}
	var user = &pb.Twitter_User{UName: temp.UName, Pwd: temp.Pwd, Posts: posts, Follow: temp.Follow}
	return user

}

func GetUserInfo(arr Tweets) []string {
	var ret []string
	for _, tweet := range arr {
		temp := tweet.User + ": " + tweet.Contents
		ret = append(ret, temp)
	}
	return ret
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Deletes(a []string, x string) []string {
	var ret []string
	for _, n := range a {
		if x != n {
			ret = append(ret, n)
		}
	}
	return ret
}
