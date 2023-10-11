package model

import "gorm.io/gorm"

type DomainSite struct {
	gorm.Model
	Site   string `gorm:"column:site;size:40;uniqueIndex;comment:站点名字"`
	Domain string `gorm:"column:domain;size:255;comment:域名"`
}

func (DomainSite) TableName() string {
	return "site_domain"
}

type DomainWarning struct {
	gorm.Model
	Domain string `gorm:"column:domain;size:255;index;comment:域名"`
	Ip     string `gorm:"column:title;size:128;comment:接口标题"`
	Path   string `gorm:"column:path;size:128;comment:接口地址"`
	Action string `gorm:"column:action;size:16;comment:请求类型"`
}

func (DomainWarning) TableName() string {
	return "domain_warning"
}
