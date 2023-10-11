package mqtt

import (
	"im/pkg/common"
	"testing"
)

func TestSendMessage(t *testing.T) {
	type args struct {
		operationID string
		mt          common.MessageType
		message     interface{}
		topic       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				operationID: "1122334455",
				mt:          common.UserInfoPush,
				message: struct {
					Name string
					Age  int
				}{},
				topic: "topic/test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMessage(tt.args.operationID, tt.args.mt, tt.args.message, tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


