package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func Test_loginHandler(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:9090/view/home.html",
		httpmock.NewStringResponder(200, `[{"username": "sneha", "password": "sneha123"}]`))

	// do stuff that makes a request to articles.json
}

func Test_logoutHandler(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:9090/login.html",
		httpmock.NewStringResponder(200, `[{"username": "sneha", "password": "sneha123"}]`))

	// do stuff that makes a request to articles.json
}
func Test_tweetHandler(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:9090/view/home.html",
		httpmock.NewStringResponder(200, `[{"username": "sneha", "password": "sneha123"}]`))

	// do stuff that makes a request to articles.json
}
