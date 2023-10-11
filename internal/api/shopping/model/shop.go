package model

import (
	"im/internal/api/user/model"
	"im/pkg/db"
)

const ShopStatusApprove = 1
const ShopStatusPass = 2
const ShopStatusDeleted = 4
const ShopStatusRejected = 3
const TeamRoleLeader = 2
const TeamRoleNobody = 1

type Shop struct {
	db.CommonModel
	Name            string     `gorm:"column:name;size:200;" json:"name"`
	Longitude       string     `gorm:"column:longitude;size:50;" json:"longitude"`
	Latitude        string     `gorm:"column:latitude;size:50;" json:"latitude"`
	Address         string     `gorm:"column:address;size:200;" json:"address"`
	License         string     `gorm:"column:license;size:500;" json:"license"`
	Image           string     `gorm:"column:image;size:2000;" json:"image"`
	Description     string     `gorm:"column:description;size:2000;" json:"description"`
	DecorationScore float64    `gorm:"column:decoration_score;" json:"decoration_score"`
	QualityScore    float64    `gorm:"column:quality_score;" json:"quality_score"`
	ServiceScore    float64    `gorm:"column:service_score;" json:"service_score"`
	Star            float64    `gorm:"column:star;" json:"star"`
	ShopType        string     `gorm:"column:;" json:"shop_type"`
	CityCode        string     `gorm:"column:city_code;" json:"city_code"`
	CreatorId       string     `gorm:"column:creator_id;" json:"creator_id"`
	CreatorUser     model.User `gorm:"foreignKey:UserID;references:CreatorId"`
	InviteCode      string     `gorm:"column:invite_code;" json:"invite_code"`
	Status          int64      `gorm:"column:status;default:1;" json:"status"`
}

func (d *Shop) TableName() string {
	return "shops"
}

type ShopTeam struct {
	db.CommonModel
	ShopID       int64      `json:"column:shop_id;default:0;" json:"shop_id"`
	Shop         Shop       `gorm:"foreignKey:ID;references:ShopID"`
	Role         int64      `gorm:"column:role;default:1;" json:"role"`
	InviteUserId string     `gorm:"column:invite_user_id;size:20;default:'';" json:"invite_user_id"`
	InviteUser   model.User `gorm:"foreignKey:UserID;references:InviteUserId"`
	UserID       string     `gorm:"column:user_id;size:100;default:'';" json:"user_id"`
	User         model.User `gorm:"foreignKey:UserID;references:UserID"`
	Status       int64      `gorm:"column:status;default:2;" json:"status"`
}

func (d *ShopTeam) TableName() string {
	return "shop_teams"
}
