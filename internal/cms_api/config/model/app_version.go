package model

import (
	"im/pkg/db"
)

type AppVersion struct {
	db.CommonModel
	Platform     int64  `gorm:"column:platform;default:2;" json:"platform"`
	Version      string `gorm:"column:version;size:20;" json:"version"`
	IsForce      int64  `gorm:"column:is_force;default:2;" json:"is_force"`
	Title        string `gorm:"column:title;size:80;" json:"title"`
	DownloadUrl  string `gorm:"column:download_url;default:'';size:255;" json:"download_url"`
	UpdateDesc   string `gorm:"column:update_desc;default:'';size:560;" json:"update_desc"`
	Status       int64  `gorm:"column:status;default:2;" json:"status"`
	DeleteStatus int64  `gorm:"column:delete_status;default:2;" json:"delete_status"`
}
