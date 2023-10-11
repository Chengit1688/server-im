package usecase

import (
	"im/config"
	"im/internal/api/chat/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"testing"
	"time"
)

func init() {
	config.Init()
	logger.Init("")
	db.Init()
}

func Test_messageUseCase_GetMsgID(t *testing.T) {
	type args struct {
		conversationType model.ConversationType
		conversationID   string
		sendTime         int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				conversationType: model.ConversationTypeSingle,
				conversationID:   "10000_10001",
				sendTime:         util.UnixMilliTime(time.Now()),
			},
		},
		{
			name: "success",
			args: args{
				conversationType: model.ConversationTypeGroup,
				conversationID:   "10001",
				sendTime:         util.UnixMilliTime(time.Now()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &messageUseCase{}
			if got := c.GetMsgID(tt.args.conversationType, tt.args.conversationID, tt.args.sendTime); got != tt.want {
				t.Errorf("GetMsgID() = %v, want %v", got, tt.want)
			}
		})
	}
}
