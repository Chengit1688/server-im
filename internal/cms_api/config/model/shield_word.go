package model

import "im/pkg/db"

type ShieldWords struct {
	db.CommonModel
	Status        int64  `gorm:"column:status;default:1;" json:"status"`
	DeleteStatus  int64  `gorm:"column:delete_status;default:2;" json:"delete_status"`
	ShieldWords   string `gorm:"column:shield_words;size:800;" json:"shield_words"`
	OperationUser string `gorm:"column:operation_user;size:80;" json:"-"`
}

func (*ShieldWords) TableName() string {
	return "shield_words"
}
