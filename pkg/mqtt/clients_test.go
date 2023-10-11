package mqtt

import (
	"reflect"
	"testing"
)

func TestGetClients(t *testing.T) {
	type args struct {
		username     string
		likeUsername string
		likeClientID string
		connState    string
		page         int
		limit        int
	}
	tests := []struct {
		name        string
		args        args
		wantClients []Client
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				username:  "xiaohu",
				connState: ConnStateTypeConnected,
			},
		},
		{
			name: "success",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClients, err := GetClients(tt.args.username, tt.args.likeUsername, tt.args.likeClientID, tt.args.connState, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClients, tt.wantClients) {
				t.Errorf("GetClients() gotClients = %v, want %v", gotClients, tt.wantClients)
			}
		})
	}
}
