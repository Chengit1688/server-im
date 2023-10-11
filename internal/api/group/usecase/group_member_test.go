package usecase

import (
	"im/config"
	"im/internal/api/group/model"
	"im/pkg/db"
	"im/pkg/logger"
	"reflect"
	"testing"
)

func init() {
	config.Init()
	logger.Init("")
	db.Init()
}

func Test_groupMemberUseCase_GetMember(t *testing.T) {
	type args struct {
		groupID  string
		memberID string
	}
	tests := []struct {
		name       string
		args       args
		wantMember *model.GroupMember
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				groupID:  "5389460160",
				memberID: "317087584766",
			},
		},
		{
			name: "failed",
			args: args{
				groupID:  "53894601601",
				memberID: "317087584766",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &groupMemberUseCase{}
			gotMember, err := c.GetMember(tt.args.groupID, tt.args.memberID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMember, tt.wantMember) {
				t.Errorf("GetMember() gotMember = %v, want %v", gotMember, tt.wantMember)
			}
		})
	}
}
