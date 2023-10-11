package model

type ClearType int

const (
	ClearTypeUser        ClearType = 1
	ClearTypeGroupMember ClearType = 2
	ClearTypeAll         ClearType = 3
)

type MultiSendRecord struct {
	ID             int    `gorm:"primary_key;autoIncrement"`
	OperateID      string `gorm:"type:varchar(25);column:operate_id"`
	Content        string `gorm:"type:Text;column:content"`
	CreatedAt      int64  `gorm:"column:created_at"`
	Username       string `gorm:"-:migration;<-:false"`
	SenderID       string `gorm:"-:migration;<-:false"`
	SenderNickname string `gorm:"-:migration;<-:false"`
}

func (MultiSendRecord) TableName() string {
	return "cms_multi_send_record"
}

type MultiSendUser struct {
	ID       int    `gorm:"primary_key;autoIncrement"`
	SenderID string `gorm:"type:varchar(25);column:sender_id"`
	RecordID int    `gorm:"type:int;column:record_id"`
}

func (MultiSendUser) TableName() string {
	return "cms_multi_send_user"
}
