package auth

import (
	"net/http"
	"twitterPt2/auth/cookie"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie.ClearSession(w)
	http.Redirect(w, r, "/login", 302)
}
