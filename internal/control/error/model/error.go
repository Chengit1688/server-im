package model

type ErrLog struct {
	Id        int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT"`
	AppName   string `gorm:"column:app_name;type:varchar(64);index"`
	UserId    string `gorm:"column:user_id;type:varchar(64);index"`
	MacType   string `gorm:"column:mac_type;type:varchar(64)"`
	PhoneType string `gorm:"column:phone_type;type:varchar(64)"`
	CreatTime int64  `gorm:"column:create_time;type:int(11)"`
	Info      string `gorm:"column:info"`
	Extra     string `gorm:"column:extra"`
}

func (ErrLog) TableName() string {
	return "errlog"
}
