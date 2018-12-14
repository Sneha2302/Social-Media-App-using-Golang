package auth

import (
	"fmt"
	"html/template"
	"net/http" 
	"twitterPt2/gRpc/client"
	"twitterPt2/auth/cookie"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method : ", r.Method)
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("view/login.html")
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
		redirectTarget := "/"
		uName := r.Form.Get("username")
			pWord := r.Form.Get("password")
			// verify credentials
			if ok := client.UserExistsRpc(uName, pWord); ok {
				fmt.Println(ok)
				cookie.SetSession(uName, w)
				redirectTarget += "home"
			}
		http.Redirect(w, r, redirectTarget, http.StatusFound)
	default:
		fmt.Fprintf(w, "method not supported.")
	}
}
