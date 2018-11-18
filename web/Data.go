package main

import (
	"fmt"
)

var userdata = make(map[string]User)

type User struct {
	username string
	password string
	follows map[string]bool
	tweets []tweet
}

type tweet struct {
	text string
}

var debugon = true //if set to true debug outputs are printed

//Function to print debug outputs if debugon=true
func debugPrint(text string){
	if(debugon){
		fmt.Println(text)
	}
}
