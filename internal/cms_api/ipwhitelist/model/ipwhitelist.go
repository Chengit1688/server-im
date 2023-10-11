package model

import (
	adminModel "im/internal/cms_api/admin/model"
)

type IPWhiteList struct {
	ID        uint             `gorm:"column:id;primaryKey;size:11" json:"id"`
	IP        string           `gorm:"column:ip;size:32;uniqueIndex:ip" json:"ip"`
	UserID    string           `gorm:"column:user_id" json:"user_id"`
	Admin     adminModel.Admin `gorm:"foreignKey:UserID;references:UserID"`
	Note      string           `gorm:"column:note;size:255" json:"note"`
	CreatedAt int64            `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64            `gorm:"column:updated_at" json:"updated_at"`
}

func (IPWhiteList) TableName() string {
	return "cms_ipwhitelist"
}
