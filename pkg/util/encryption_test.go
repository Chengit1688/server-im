package util

import (
	"im/pkg/common"
	"testing"
)

func TestGetPassword(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				password: "Qewr@1234",
				salt:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := GetPassword(tt.args.password, tt.args.salt)
			t.Errorf("GetPassword() password_encode = %v", pass)
		})
	}
}

func TestCheckPassword(t *testing.T) {
	type args struct {
		dbPassword   string
		userPassword string
		salt         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				dbPassword:   "412d81cb98ebc6564ea8876df9f23324",
				salt:         "BhsbysoP",
				userPassword: "Qewr@1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CheckPassword(tt.args.dbPassword, tt.args.userPassword, tt.args.salt)
			t.Errorf("TestCheckPassword() resule = %v", res)
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		data []byte
		key  []byte
	}
	tests := []struct {
		name        string
		args        args
		wantContent string
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				data: []byte("你好123abc"),
				key:  common.ContentKey,
			},
			wantContent: "fOywEkN/Tj0gSGn+H8jVuA==",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := Encrypt(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContent != tt.wantContent {
				t.Errorf("Encrypt() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		dataStr string
		key     []byte
	}
	tests := []struct {
		name        string
		args        args
		wantContent string
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				dataStr: "fOywEkN/Tj0gSGn+H8jVuA==",
				key:     common.ContentKey,
			},
			wantContent: "你好123abc",
			wantErr:     false,
		},
		{
			name: "success",
			args: args{
				dataStr: "GcJBsoFX0q0VxP6+Pn2tUt7egy3DGwmzIXPIVopZbxO5qiS0SSoTQNa51d4qNApZ/HC112Z0y4XdvvsGZxR8/qnowkkgNIG1dSI3X5OaUMVgc9lGhNAl+UfFxbxh1FIvaziHFFMOhyQu+dRZF53pIyyEWMibqYKOliV8U5WElemdeWiE5WD9tXQ+WyyDm8Mlttu6HFdiNCW5YA9Ttrs48Q==",
				key:     common.ContentKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := Decrypt(tt.args.dataStr, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContent != tt.wantContent {
				t.Errorf("Decrypt() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}
