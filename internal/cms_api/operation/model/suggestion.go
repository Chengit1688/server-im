package model

import (
	"im/pkg/db"
)

func (Suggestion) TableName() string {
	return "suggestions"
}

type Suggestion struct {
	db.CommonModel
	UserID     string `gorm:"column:user_id;size:80;default:'';" json:"user_id"`
	Account    string `gorm:"column:account;size:80;default:'';" json:"account"`
	NickName   string `gorm:"column:nick_name;size:255;default:'';" json:"nick_name"`
	Content    string `gorm:"column:content;size:500;default:'';" json:"content"`
	Brand      string `gorm:"column:brand;size:500;default:'';" json:"brand"`
	Platform   int64  `gorm:"column:platform;" json:"platform"`
	AppVersion string `gorm:"column:app_version;size:50;default:'';" json:"version"`
}
