package model

import "gorm.io/gorm"

type CmsMenu struct {
	gorm.Model
	Name       string      `gorm:"size:128;" json:"name"`
	Title      string      `gorm:"size:128;" json:"title"`
	Icon       string      `gorm:"size:128;" json:"icon"`
	Path       string      `gorm:"size:128;" json:"path"`
	Paths      string      `gorm:"size:128;" json:"paths"`
	Type       int         `gorm:"size:1;" json:"type"`
	Action     string      `gorm:"size:16;" json:"action"`
	Permission string      `gorm:"size:255;" json:"permission"`
	ParentId   int         `gorm:"size:11;" json:"parent_id"`
	NoCache    int         `gorm:"size:1;" json:"no_cache"`
	Component  string      `gorm:"size:255;" json:"component"`
	Sort       int         `gorm:"size:4;" json:"sort"`
	Visible    int         `gorm:"size:1;DEFAULT:1;" json:"visible"`
	Hidden     int         `gorm:"size:1;DEFAULT:1;" json:"hidden"`
	IsFrame    int         `gorm:"size:1;DEFAULT:1;" json:"is_frame"`
	Apis       *[]AdminApi `gorm:"many2many:cms_menu_api;foreignKey:ID;joinForeignKey:menu_id;references:ID;joinReferences:api_id;" json:"apis"`
}

func (CmsMenu) TableName() string {
	return "cms_menus"
}

type AdminApi struct {
	gorm.Model
	Handle string `gorm:"column:handle;size:128;comment:接口函数"`
	Title  string `gorm:"column:title;size:128;comment:接口标题"`
	Path   string `gorm:"column:path;size:128;comment:接口地址"`
	Action string `gorm:"column:action;size:16;comment:请求类型"`
}

func (AdminApi) TableName() string {
	return "cms_admin_apis"
}

type Config struct {
	ID         uint   `gorm:"column:id;primaryKey;size:11"`
	Name       string `gorm:"column:name;size:64;uniqueIndex:name"`
	Value      string `gorm:"column:value;size:255"`
	Note       string `gorm:"column:note;size:255"`
	UpdateTime int64  `gorm:"column:update_time"`
}

func (Config) TableName() string {
	return "config"
}

type CmsMenuApi struct {
	ID     int `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT"`
	MenuID int `gorm:"column:menu_id;type:int(11)"`
	ApiID  int `gorm:"column:api_id;type:int(11)"`
}

func (m *CmsMenuApi) TableName() string {
	return "cms_menu_api"
}
