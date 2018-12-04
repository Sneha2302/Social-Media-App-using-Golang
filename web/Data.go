package main

import (
	"fmt"
)

var debugon = true //if set to true debug outputs are printed

//Function to print debug outputs if debugon=true
func debugPrint(text string) {
	if debugon {
		fmt.Println(text)
	}
}
