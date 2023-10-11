package model

import (
	"im/pkg/db"
)

type InviteCode struct {
	db.CommonModel
	InviteCode     string `gorm:"column:invite_code;size:80;uniqueIndex;" json:"invite_code"`
	DefaultGroups  string `gorm:"column:default_groups;size:720;" json:"default_groups"`
	DefaultFriends string `gorm:"column:default_friends;size:720;" json:"default_friends"`
	GreetMsg       string `gorm:"column:greet_msg;size:255;" json:"greet_msg"`
	Remarks        string `gorm:"column:remarks;size:500;default:'';size:320;" json:"remarks"`
	Status         int64  `gorm:"column:status;default:1;" json:"status"`
	DeleteStatus   int64  `gorm:"column:delete_status;default:2;" json:"delete_status"`
	OperationUser  string `gorm:"column:operation_user;size:80;" json:"operation_user"`
	IsOpenTurn     int64  `gorm:"column:is_open_turn;default:2;" json:"is_open_turn"`
	FriendIndex    int64  `gorm:"column:friend_index;default:0;" json:"friend_index"`
	GroupIndex     int64  `gorm:"column:group_index;default:0;" json:"group_index"`
}

func (InviteCode) TableName() string {
	return "invite_codes"
}
