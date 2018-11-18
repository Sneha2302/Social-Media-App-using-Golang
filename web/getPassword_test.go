package main

import (
	"fmt"
	"testing"
)

func Test_getPassword(t *testing.T) {
	type args struct {
		usrname string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{

		{
			"sneha",
			args{"sneha"},
			true,
			"sneha12",
		},
		{
			"mayank",
			args{"mayank"},
			true,
			"1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getPassword(tt.args.usrname)
			if got == tt.want {
				fmt.Printf("Test Passed : got = %v, want %v", got, tt.want)
			}
			if got != tt.want {
				t.Errorf(" Test Failed :getPassword() got = %v, want %v", got, tt.want)
			}
			if got1 == tt.want1 {
				fmt.Printf("Test Passed : got1 = %v, want %v", got1, tt.want1)
			}

			if got1 != tt.want1 {
				t.Errorf("Test failed: getPassword() got1 = %v, want %v", got1, tt.want1)
			}

		})
	}
}
