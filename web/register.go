package web

import (
	"fmt"
	"html/template"
	"net/http"
	"social_media_app-golang/gRpc/client"
)

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method : ", r.Method)
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("view/register.html")
		if err != nil {
			fmt.Fprintf(w, "Error : %v\n", err)
			return
		}
		t.Execute(w, nil)
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Error: %v\n", err)
			return
		}
		redirectTarget := "/login"
		uName := r.Form.Get("username")
		pWord1 := r.Form.Get("password1")
		pWord2 := r.Form.Get("password2")
		if ok := client.RegisterRpc(uName, pWord1, pWord2); ok {
			fmt.Println("User registered success")
		} else {
			fmt.Println("User registration failed")
			fmt.Println(ok)
		}
		http.Redirect(w, r, redirectTarget, http.StatusFound)
	default:
		fmt.Fprintf(w, "method not supported.")
	}
}
