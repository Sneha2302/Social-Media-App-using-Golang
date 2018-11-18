package main

import "testing"

func Test_deleteUser(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			"adsv",
			args{"sneha"},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteUser(tt.args.username); got != tt.want {
				t.Errorf("deleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
