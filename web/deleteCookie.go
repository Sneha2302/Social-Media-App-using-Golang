package main

import (
	"net/http"
)

func deleteCookie(w http.ResponseWriter){
	cookie := http.Cookie{Name: "username", MaxAge: -1}
	http.SetCookie(w, &cookie)
	debugPrint("Debug:Cookie Deleted")
	return
}