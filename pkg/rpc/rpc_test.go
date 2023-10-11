package rpc

import (
	"im/pkg/common"
	"testing"
)

func TestSendMessageToConnectServer(t *testing.T) {
	type args struct {
		operationID string
		mt          common.MessageType
		message     interface{}
		userIDs     []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				operationID: "112233",
				mt:          common.UserInfoPush,
				userIDs:     []string{"001", "002"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{}
			message.Username = "hello"
			message.Password = "hello123"
			tt.args.message = message

			SendMessageToConnectServer(tt.args.operationID, tt.args.mt, tt.args.message, tt.args.userIDs...)
		})
	}
}
