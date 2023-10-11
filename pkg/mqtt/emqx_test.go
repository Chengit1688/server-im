package mqtt

import (
	"im/config"
	"im/pkg/db"
	"im/pkg/logger"
	"testing"
)

func init() {
	config.Init()
	logger.Init("")
	db.Init()
}

func TestConnect(t *testing.T) {
	type args struct {
		userID   string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				clientID: "xiaohu0002",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Connect(tt.args.userID, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	type args struct {
		topic    string
		qos      int
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				topic:    "topic_xiaohu001",
				qos:      2,
				clientID: "xiaohu0002",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Subscribe(tt.args.topic, tt.args.qos, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublish(t *testing.T) {
	type args struct {
		topic    string
		qos      int
		clientID string
		payload  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				topic:    "topic_xiaohu001",
				qos:      2,
				clientID: "system",
				payload:  "hello xiaohu",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Publish(tt.args.topic, tt.args.qos, tt.args.clientID, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
