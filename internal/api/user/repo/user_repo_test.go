package repo

import (
	"im/config"
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

func Test_userRepo_ListUserID(t *testing.T) {
	tests := []struct {
		name           string
		wantUserIDList []string
		wantErr        bool
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{}
			gotUserIDList, err := r.ListUserID()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserIDList, tt.wantUserIDList) {
				t.Errorf("ListUserID() gotUserIDList = %v, want %v", gotUserIDList, tt.wantUserIDList)
			}
		})
	}
}
