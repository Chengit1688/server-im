package model

import (
	"gorm.io/gorm"
)

type Discover struct {
	gorm.Model
	Name string `gorm:"column:name;size:50;" json:"name"`
	Url  string `gorm:"column:url;size:255" json:"url"`
	Icon string `gorm:"column:icon;size:255" json:"icon"`
	Sort int    `gorm:"column:sort;" json:"sort"`
}

func (Discover) TableName() string {
	return "discover"
}

type PrizeList struct {
	ID              int64  `gorm:"column:id;primarykey" json:"id"`
	CreatedAt       int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       int64  `gorm:"column:updated_at" json:"updated_at"`
	Name            string `gorm:"column:name;default:'';size:100;" json:"name"`
	FictitiousValue string `gorm:"column:fictitious_value;default:'';size:100;" json:"fictitious_value"`
	Cost            int64  `gorm:"column:cost;default:0;" json:"cost"`
	Icon            string `gorm:"column:icon;default:'';size:255;" json:"icon"`
	Describe        string `gorm:"column:describe;default:'';size:255;"  json:"describe"`
	IsFictitious    int    `gorm:"column:is_fictitious;default:1;" json:"is_fictitious"`
	Status          int64  `gorm:"column:status;default:1;" json:"status"`
}

func (PrizeList) TableName() string {
	return "prize_lists"
}

type RedeemPrizeLog struct {
	ID            int64     `gorm:"column:id;primarykey" json:"id"`
	CreatedAt     int64     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     int64     `gorm:"column:updated_at" json:"updated_at"`
	Name          string    `gorm:"column:name;default:'';size:100;" json:"name"`
	Icon          string    `gorm:"column:icon;default:'';size:255;" json:"icon"`
	IsFictitious  int       `gorm:"column:is_fictitious;default:1;" json:"is_fictitious"`
	UserId        string    `gorm:"column:user_id;" json:"user_id"`
	UserName      string    `gorm:"column:user_name;size:100;" json:"user_name"`
	Address       string    `gorm:"column:address;size:500;" json:"address"`
	Mobile        string    `gorm:"column:mobile;size:50;" json:"mobile"`
	PrizeId       int64     `gorm:"column:prize_id;default:0;" json:"prize_id"`
	Prize         PrizeList `gorm:"foreignKey:ID;references:PrizeId" json:"-"`
	Cost          int64     `gorm:"column:cost;default:0;" json:"cost"`
	ExpressNumber string    `gorm:"column:express_number;default:'';size:80;" json:"express_number"`
	Status        int64     `gorm:"column:status;default:1;" json:"status"`
}

func (RedeemPrizeLog) TableName() string {
	return "redeem_prize_logs"
}
