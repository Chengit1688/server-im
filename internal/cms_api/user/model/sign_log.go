package model

import "im/pkg/db"

type SignLog struct {
	db.CommonModel
	Time   string `gorm:"column:time;size:80;" json:"time"`
	Year   int    `gorm:"column:year;" json:"year"`
	Month  int    `gorm:"column:month;" json:"month"`
	Day    int    `gorm:"column:day;" json:"day"`
	UserId string `gorm:"column:user_id;" json:"user_id"`
	Reward int64  `gorm:"column:reward;default:1;" json:"reward"`
	Status string `gorm:"column:status;default:1;" json:"status"`
}

func (*SignLog) TableName() string {
	return "sign_logs"
}
