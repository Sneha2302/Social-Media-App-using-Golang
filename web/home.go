package web

import (
	"twitterPt2/auth/cookie"
	"fmt"
	"html/template"
	"net/http"
	"twitterPt2/gRpc/client"
	"twitterPt2/storage"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	uName := cookie.GetUserName(r)
	if uName != "" {
		userInfo := client.GetUserRpc(uName)
		fmt.Println("method:::::", r.Method)

		switch r.Method {
		case "GET":
			t, err := template.ParseFiles("view/home.html")
			if err != nil {
				fmt.Fprintf(w, "Error : %v\n", err)
				return
			}
			userInfo := client.GetUserInfoRpc(uName)
			t.Execute(w, userInfo)
		case "POST":
			r.ParseForm()
			submitType := r.Form.Get("submit")
			fmt.Println(submitType)
			redirectUrl := r.URL.Path
			switch submitType {
			case "follow":
				person := r.Form.Get("unfollow")
				client.FollowUserRpc(uName, person)
			case "tweet":
				var userTweet = storage.TwitterPosts{}
				userTweet.Contents = r.Form.Get("tweet_text")
				fmt.Println("Inside tweet post")
				if userTweet.Contents != "" {
					userTweet.Date = time.Now().Unix()
					userTweet.User = uName
					userInfo.Post = append(userInfo.Post, userTweet) 
					client.TweetRpc(uName, userInfo)
					fmt.Println("Posts", userInfo.Post)
				}
				redirectUrl = "/home"
			}
			http.Redirect(w, r, redirectUrl, 302)
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}
