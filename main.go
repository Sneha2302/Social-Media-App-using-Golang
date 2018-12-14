package main

import (
	"twitterPt2/auth"
	"net/http"
	"twitterPt2/web"
)

func main() {
	http.HandleFunc("/", auth.Login)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/register", web.Register)
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/profile", web.Profile)
	http.HandleFunc("/logout", auth.Logout)
	http.ListenAndServe(":9090", nil)
}
