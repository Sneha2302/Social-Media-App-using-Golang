package main

import (
	"fmt"
	"testing"
)

func Test_addUser(t *testing.T) {
	type args struct {
		usrname string
		pwd     string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"adsv",
			args{"sneha", "sneha12"},
			1,
		},

		{
			"adsv11",
			args{"shradha", "123"},
			1,
		},
		{
			"adsv12",
			args{"mayank", "1234"},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addUser(tt.args.usrname, tt.args.pwd); got != tt.want {

				fmt.Printf("Test pass yay!!")
				t.Errorf("addUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
