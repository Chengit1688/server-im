package util

import (
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	type args struct {
		userID   string
		duration time.Duration
	}
	tests := []struct {
		name            string
		args            args
		wantTokenString string
		wantErr         bool
	}{
		{
			name: "success",
			args: args{
				userID:   "001",
				duration: time.Hour * 24 * 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenString, err := CreateToken(tt.args.userID, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTokenString != tt.wantTokenString {
				t.Errorf("CreateToken() gotTokenString = %v, want %v", gotTokenString, tt.wantTokenString)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIwMDEiLCJleHAiOjE2NzI3NDE5MTZ9.izVrWKaTtCf_oVqCxLKE0Lx74HQORYgU1oYv5-FJEY8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseToken(tt.args.tokenString); (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
