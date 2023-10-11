package model

import (
	"im/pkg/db"
)

type DefaultFriend struct {
	db.CommonModel
	UserPrimaryId   int64  `gorm:"column:user_primary_id;uniqueIndex;" json:"user_primary_id"`
	UserId          string `gorm:"column:user_id;size:80;uniqueIndex;" json:"user_id"`
	GreetMsg        string `gorm:"column:greet_msg;" json:"greet_msg"`
	OperationUserId string `gorm:"column:operation_user_id;" json:"operation_user_id"`
	Remarks         string `gorm:"column:remarks;default:'';size:320;" json:"remarks"`
}

func (d *DefaultFriend) TableName() string {
	return "default_friends"
}
