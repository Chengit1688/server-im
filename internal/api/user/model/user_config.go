package model

import "im/pkg/db"

type UserConfig struct {
	db.CommonModel
	UserId  string `gorm:"column:user_id;uniqueIndex;size:50"`
	Content string `gorm:"column:content;type:text"`
	Version int64  `gorm:"column:version;default:1"`
}
