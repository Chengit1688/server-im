package model

type Face2faceGroup struct {
	Id           int64  `gorm:"column:id;autoIncrement;uniqueIndex;" json:"id"`
	TmpGroupId   string `gorm:"column:tmp_group_id;size:20;uniqueIndex;" json:"tmp_group_id"`
	Name         string `gorm:"column:name;type:varchar(255);index" json:"name"`
	UserId       string `gorm:"column:user_id;type:varchar(50);" json:"user_id"`
	CreateUserId string `gorm:"column:create_user_id;type:varchar(50);" json:"create_user_id"`
	Longitude    string `gorm:"column:longitude;size:50;" json:"longitude"`
	Latitude     string `gorm:"column:latitude;size:50;" json:"latitude"`
	RandomNumber string `gorm:"column:random_number;size:20;" json:"random_number"`
	CreatedAt    int64  `gorm:"column:created_at" json:"created_at"`
	ExpireTime   int64  `gorm:"column:expire_time" json:"expire_time"`
	UpdatedAt    int64  `gorm:"column:updated_at" json:"updated_at"`
	Status       int    `gorm:"column:status;default:1;type:tinyint(3);" json:"status"`
}

func (d *Face2faceGroup) TableName() string {
	return "face2face_group"
}
