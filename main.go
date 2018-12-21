package main

import (
	"net/http"
	"social_media_app-golang/auth"
	"social_media_app-golang/web"
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
