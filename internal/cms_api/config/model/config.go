package model

type Config struct {
	ID         uint   `gorm:"column:id;primaryKey;size:11"`
	Name       string `gorm:"column:name;size:64;uniqueIndex:name"`
	Value      string `gorm:"column:value;type:longtext"`
	Note       string `gorm:"column:note;size:255"`
	UpdateTime int64  `gorm:"column:update_time"`
}

func (Config) TableName() string {
	return "cms_config"
}

type SettingConfig struct {
	ID         uint   `gorm:"column:id;primaryKey;size:11"`
	ConfigType string `gorm:"column:config_type;size:255"`
	Content    string `gorm:"column:content;type:text"`
}
