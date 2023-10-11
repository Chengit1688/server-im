package model

type Link struct {
	MsgID            string           `gorm:"size:64;index:idx_msg_id,unique" json:"msg_id"`
	ClientMsgID      string           `gorm:"size:64;" json:"client_msg_id"`
	ConversationType ConversationType `gorm:"type:tinyint;index:idx_conversation_type_conversation_id_seq" json:"conversation_type"`
	ConversationID   string           `gorm:"size:40;index:idx_conversation_type_conversation_id_seq" json:"conversation_id"`
	SendID           string           `gorm:"size:20" json:"send_id"`
	SendTime         int64            `gorm:"type:bigint;index:idx_send_time" json:"send_time"`
	Type             MessageType      `gorm:"type:smallint" json:"type"`

	Title string `gorm:"type:varchar(255);default:''" json:"title"`

	Link string `gorm:"type:text;not null"  json:"Link"`

	Cover string `gorm:"type:text"  json:"cover"`

	Desc string `gorm:"type:text"  json:"desc"`

	Favicon string `gorm:"type:text"  json:"favicon"`

	Keyword string `gorm:"type:text"  json:"keyword"`

	Status int `gorm:"type:tinyint(1);default:0"  json:"status"`

	Content string `gorm:"type:text"  json:"content"`
}
