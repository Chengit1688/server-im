package model

import (
	"im/pkg/db"
)

type SettingConfig struct {
	db.CommonModel
	Content    string `gorm:"column:content;type:text;" json:"content"`
	ConfigType string `gorm:"column:config_type;uniqueIndex;size:50;" json:"config_type"`
}
