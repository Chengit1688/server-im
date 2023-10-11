package model

type IPBlackList struct {
	ID   uint   `gorm:"column:id;primaryKey;size:11" json:"id"`
	IP   string `gorm:"column:ip;size:32;uniqueIndex:ip" json:"ip"`
	Note string `gorm:"column:note;size:255" json:"note"`
}

func (IPBlackList) TableName() string {
	return "cms_ipblacklist"
}
