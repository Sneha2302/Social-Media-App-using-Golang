package main

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	// "strings"
	// "os"
	// "time"
	// "golang.org/x/net/context"
	// "google.golang.org/grpc"
	//pb "project/utils/ProtoDef"
)
const (
 	address     = "localhost:50051"
 	defaultName = "world"
 )

 // var rpcCaller pb.GreeterClient

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//User has directly come to login page
		t, _ := template.ParseFiles("../view/login.html")
		t.Execute(w, nil)
	} else {
		//User has come to login via Post
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		usr := r.Form["username"][0]
		pwd := r.Form["password"][0]
		ok, actualPassword := getPassword(usr)

		//User does not exist - send to register page
		if !ok {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if pwd == actualPassword {
			//Login successful, set cookie and goto tweet page 
			cookie := http.Cookie{Name: "username", Value: usr, MaxAge: 2800}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/tweet", http.StatusSeeOther)
			return
		} else {
			//Login unsuccessful go back to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		// reply, err := rpcCaller.Login(ctx, &pb.Credentials{Uname: usr, Pwd: pwd})
		// //User does not exist - send to registration page
		// if err != nil {
		// 	fmt.Println("Debug: Login rpc failed", err.Error())
		// 	if err.Error() == "Wrong Password" {
		// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 		return
		// 	} else {
		// 		http.Redirect(w, r, "/registration", http.StatusSeeOther)
		// 		return
		// 	}

		// } else if err == nil && reply.Status == true {
		// 	debugPrint("debug: user successfully logged in")
		// 	expiration := 2800
		// 	cookie := http.Cookie{Name: "username", Value: usr, MaxAge: expiration}
		// 	http.SetCookie(w, &cookie)
		// 	http.Redirect(w, r, "/tweet", http.StatusSeeOther)
		// 	return
		// } else {
		// 	log.Println("Major issue")
		// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 	return
		// }
	}
}

//Handler for register
func registerHandler(w http.ResponseWriter, r *http.Request) {
	//Get request method, type of request
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("../view/register.html")
		t.Execute(w, nil)
		return
	} else {
		r.ParseForm()
		if debugon {
			fmt.Println("username in post: ", r.Form["username"][0])
			fmt.Println("password in post: ", r.Form["password_1"][0])
		}

		//Check for non-empty username and password values
		if len(r.Form["username"][0]) == 0 || len(r.Form["password_1"][0]) == 0 {
			if debugon {
				fmt.Println("Empty Username or Password value")
			}
			//TODO: Remove Alert?
			//fmt.Fprintln(w, "<script>alert(\"Please enter a valid Username and Password\")</script>")
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		//Adding username and password to the Map
		result := addUser(r.Form["username"][0], r.Form["password_1"][0])
		if result == 1 {
			//Successfully added user to the map, redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			//User already exists send back to register page
			if debugon {
				fmt.Println("The user already exists")
			}
			//TODO: Remove Alert?
			//fmt.Fprintln(w, "<script>alert(\"User already exists\")</script>")
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
	}
}

//Home Page handler
func tweetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	cookie, ok := r.Cookie("username")
	fmt.Println(w)
	fmt.Println(r)

	//Cookie does not exist re-direct to login
	if (ok != nil) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

		// Get User Data from Map
	username := cookie.Value
	user, isUserPresent := userdata[username]

		// If map returns false, the account has been deleted. Redirect to register.
	if !isUserPresent {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		//Re-direct to homepage from login. Display Home
		t, _ := template.ParseFiles("../view/home.html")
		t.Execute(w, nil)
	} else {
		//Post: submission of new tweet. Save the tweet and then display Home.
		r.ParseForm()
		newTweet := tweet{text: r.Form["tweet"][0]}
		user.tweets = append(user.tweets, newTweet)
		userdata[username] = user
		t, _ := template.ParseFiles("../view/home.html")
		t.Execute(w, nil)
		fmt.Println(user.tweets)
	}

	fmt.Fprint(w, "<div class='login'>")
	//Display all the tweets
	if len(user.tweets) != 0 {
		fmt.Fprint(w, "<h>Tweets by you: <h><br />")
	}else{
		fmt.Fprint(w, "<h>Tweet Something? <h><br />")
	}
	for _, dispTweet := range user.tweets {
		fmt.Fprint(w, dispTweet.text)
		fmt.Fprint(w, "<br />")
	}
	fmt.Fprint(w, "<br/><br />")

	followingusers := false
	for _, v := range user.follows {
		if (v == true) {
			followingusers = true
			break
		}
	}
		fmt.Println(followingusers)
	if followingusers {
		fmt.Fprint(w, "<h>Tweets by friends: <h><br/>")
		fmt.Println(user.follows)
		for friend, _ := range user.follows {

			fuser, present := userdata[friend]
				fmt.Println(friend)
				fmt.Println(present)
				fmt.Println(userdata[friend])
			if (user.follows[friend] == true) {
				fmt.Fprint(w, "<br/>"+fuser.username+":"+"<br/>")
				for _, dispftweet := range fuser.tweets {
					fmt.Fprint(w, dispftweet.text)
					fmt.Fprint(w, "<br />")
				}
			}
		}
	}else{
		fmt.Fprint(w, "<h><a href='/users'> Users </a><h><br />")
	}

	fmt.Fprint(w, "</div>")
}

//Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//Print request method
	fmt.Println("method:", r.Method)
	deleteCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	//Print request method
	fmt.Println("method:", r.Method)
	cookie, ok := r.Cookie("username")
	if ok != nil {
		//Cookie does not exist re-direct to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	user := userdata[username]
	t, _ := template.ParseFiles("../view/Home.html")
	t.Execute(w, nil)
	follow := r.URL.Query().Get("follow")
	unfollow := r.URL.Query().Get("unfollow")
	fmt.Println("follow" + follow)
	if follow != "" {
		user.follows[follow] = true
		fmt.Println("Added" + follow)
	}
	if  unfollow != "" {
		user.follows[unfollow] = false
		fmt.Println("Removed " + unfollow)
		fmt.Println(user.follows)
	}
	fmt.Fprint(w, "<div class='login'>")
	fmt.Fprint(w, "<h>Followed Users: <h><br/>")
	for k, _ := range userdata {
		_, ok := user.follows[k]
		fmt.Println(k)
		fmt.Println(user.follows[k])
		fmt.Println(ok)
		if (user.follows[k] == false && k != user.username) {
			fmt.Println("we be in false")
			//user is not already following the person. Checking if the user1 from userdata exists in current users follows
			//fmt.Fprint(w, k)
			fmt.Fprintf(w, "%s <a href=users?follow=%s&unfollow=>Follow</a>", k, k)
			fmt.Fprint(w, "</br>")
		} else if (user.follows[k] == true && k != user.username) {
			fmt.Println("we be in true")
			fmt.Fprintf(w, "%s <a href=users?unfollow=%s&follow=>Unfollow</a>", k, k)
			fmt.Fprint(w, "</br>")
		}
	}
	fmt.Fprint(w, "</div>")
}

//Delete Account handler
func deleteAccHandler(w http.ResponseWriter, r *http.Request) {
	//get request method
	fmt.Println("Method:", r.Method)

	//Get cookie to identify the user
	cookie, ok := r.Cookie("username")
	if ok != nil {
		//Cookie does not exist, re-direct to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	//Remove user from the user Map
	deleteUser(username)

	//Delete cookie and redirect to register
	deleteCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func main() {

	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// //c := pb.NewGreeterClient(conn)
	// rpcCaller = pb.NewGreeterClient(conn)

	// // Contact the server and print out its response. TO test if RPC is working
	// name := defaultName
	// if len(os.Args) > 1 {
	// 	name = os.Args[1]
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// r, err := rpcCaller.SayHello(ctx, &pb.HelloRequest{Name: name})
	// if err != nil {
	// 	log.Fatalf("Please bring up the Back-End Server first. Could not greet: %v", err)
	// }
	// log.Printf("RPC is working %s", r.Message)
	//end of test RPC

	// All handler functions
	http.HandleFunc("/tweet", tweetHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/deleteAcc", deleteAccHandler)

	//Our server listens on this port
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
