package auth

import (
	"net/http"
	"social_media_app-golang/auth/cookie"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie.ClearSession(w)
	http.Redirect(w, r, "/login", 302)
}
