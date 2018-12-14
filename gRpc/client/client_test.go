package client

import (
	"testing"
	// "fmt"
	"strconv"
	"sync"
	"twitterPt2/storage"
)



var concurrencyFactor = 100

var emptyStr = ""
var username1 = "Sneha"
var username2 = "Shradha"
var username3 = "KS"
var username4 = "SK"

var password1 = "100"
var password2 = "200"
var password3 = "3"
var password4 = "4"

var post1 = storage.TwitterPosts{Contents: "hi", Date: 1, User: "Sneha"}
var post2 = storage.TwitterPosts{Contents: "hello", Date: 2, User: "Shradha"}
var postList = storage.TwitterPosts{post1, post2}
var user = storage.Twitter_User{UName: username1, Pwd: password1, Post: postList, Follow: []string{username1, username2}}

/*
concurrency test:
go run() {}

t.run():
*/

func TestAddUser(t *testing.T) {
	t.Run("EdgeCaseTest", func(t *testing.T) {
		if result := RegisterRpc(username1, password1, password1); result != true {
			t.Errorf("User Added successfully")
		}
		if result := RegisterRpc(username2, password2, password2); result != true {
			t.Errorf("User Added successfully ")
		}
		if result := RegisterRpc(username1, password1, password1); result != false {
			t.Errorf("User Added successfully")
		}
		if result := RegisterRpc(emptyStr, password2, password1); result != false {
			t.Errorf("User Added successfully")
		}
		if result := RegisterRpc(username2, emptyStr, password1); result != false {
			t.Errorf("No password found: empty")
		}
		if result := RegisterRpc(username4, password4, password4); result != true {
			t.Errorf("Test Passed")
		}
	})
	t.Run("ConcurrencyTest", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrencyFactor)
		for i := 0; i < concurrencyFactor; i++ {
			var testName = username1 + strconv.Itoa(i)
			var testPassword = password1 + strconv.Itoa(i)
			go func(name string, password string) {
				defer wg.Done()
				RegisterRpc(name, password, password)
			}(testName, testPassword)
		}
		wg.Wait()
	})
}

func TestUserExists(t *testing.T) {
	t.Run("TestCase", func(t *testing.T) {
		if result := UserExistsRpc(username1, password1); result != true {
			t.Errorf("User exists")
		}
		if result := UserExistsRpc(username1, password2); result != false {
			t.Errorf("User Exists")
		}
		if result := UserExistsRpc(emptyStr, password1); result != false {
			t.Errorf("User Exists")
		}
	})
	t.Run("ConcurrencyTest", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrencyFactor)
		for i := 0; i < concurrencyFactor; i++ {
			var testName = username1 + strconv.Itoa(i)
			var expect = true
			var testPassword = password1 + strconv.Itoa(i)
			if i > concurrencyFactor/2 {
				expect = false
				testPassword = testPassword + strconv.Itoa(i)
			}
			go func(name string, password string, expect bool) {
				defer wg.Done()
				if result := RpcHasUser(name, password); result != expect {
					t.Errorf("User exists /not exists")
				}
			}(testName, testPassword, expect)
		}
		wg.Wait()
	})
}

func TestTweet(t *testing.T) {
	t.Run("NormalCaseTest", func(t *testing.T) {
		if result := TweetRpc(username1, user); result != true {
			t.Errorf("Should update correctly")
		}
	})
	t.Run("ConcurrencyTest", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrencyFactor)
		for i := 0; i < concurrencyFactor; i++ {
			go func(name string, user storage.User) {
				defer wg.Done()
				if result := TweetRpc(username1, user); result != true {
					t.Errorf("Should update correctly")
				}
			}(username1, user)
		}
		wg.Wait()
	})
}

func TestFollow(t *testing.T) {
	t.Run("NormalCaseTest", func(t *testing.T) {
		if result := FollowUserRpc(username1, username2); result != false {
			t.Errorf("Expected Follwer already followed!")
		}
		if result := FollowUserRpc(username1, username3); result != true {
			t.Errorf("Expected Followed successfully")
		}
		if result := FollowUserRpc(username2, username4); result != true {
			t.Errorf("Expected Followed successfully!")
		}
	})
	t.Run("ConcurrencyTest", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrencyFactor)
		for i := 0; i < concurrencyFactor; i++ {
			go func(username1 string, username2 string) {
				defer wg.Done()
				if result := FollowUserRpc(username1, username2); result != false {
					t.Errorf("User already being followed")
				}
			}(username1, username2)
		}
		wg.Wait()
	})
}

func TestUnFollow(t *testing.T) {
	t.Run("NormalCaseTest", func(t *testing.T) {
		if result := FollowUserRpc(username1, username2); result != true {
			t.Errorf("Test successful!")
		}
		if result := FollowUserRpc(username3, username4); result != false {
			t.Errorf("Not following !")
		}
	})
	t.Run("ConcurrencyTest", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrencyFactor)
		for i := 0; i < concurrencyFactor; i++ {
			go func(username1 string, username2 string) {
				defer wg.Done()
				if result := FollowUserRpc(username1, username2); result != false {
					t.Errorf("Follwer already followed!")
				}
			}(username1, username2)
		}
		wg.Wait()
	})
}

func TestGetUserInfo(t *testing.T) {
	t.Run("Test", func(t *testing.T) {
		if pg := GetUserInfoRpc(username1); len(pg.Post) != 2 {
			t.Errorf("Correct Pages loaded !")
		}
	})
	t.Run("ConcurrencyTest", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(concurrencyFactor)
		for i := 0; i < concurrencyFactor; i++ {
			go func(username1 string) {
				defer wg.Done()
				func GetUserInfoRpc(uName string) storage.Pages {
				if pg := (username1); len(pg.Post) != 2 {
					t.Errorf("Correct user info page not loaded !")
				}
			}(username1)
		}
		wg.Wait()
	})
}
