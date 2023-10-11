package model

import (
	"im/internal/api/user/model"
	"im/pkg/db"
)

const StatusOn = 1
const StatusOff = 2

type News struct {
	db.CommonModel
	CreateUserID string     `gorm:"column:create_user_id;size:50;index:idx_create_user_id;" json:"create_user_id"`
	CreatorUser  model.User `gorm:"foreignKey:UserID;references:CreateUserID"`
	DeleteUserID string     `gorm:"column:delete_user_id;size:50;index:idx_delete_user_id;" json:"delete_user_id"`
	Title        string     `gorm:"column:title;size:200;index:idx_title;" json:"title"`
	Content      string     `gorm:"column:content;type:text;" json:"content"`
	ViewTotal    int64      `gorm:"column:view_total;type:bigint;" json:"view_total"`
	CategoryID   int64      `gorm:"column:category_id;index:idx_category_id;" json:"category_id"`
	Image        string     `gorm:"column:image;size:1500;default:'';" json:"image"`
	Video        string     `gorm:"column:video;size:1000;default:'';" json:"video"`
	Status       int64      `gorm:"column:status;" json:"status"`
}

func (d *News) TableName() string {
	return "news"
}
