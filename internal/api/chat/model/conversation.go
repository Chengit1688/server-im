package model

type ConversationType int

const (
	ConversationTypeNone    ConversationType = 0
	ConversationTypeSingle  ConversationType = 1
	ConversationTypeGroup   ConversationType = 2
	ConversationTypeChannel ConversationType = 3
)

type Conversation struct {
	ID               uint             `gorm:"primaryKey" json:"-"`
	DeletedAt        int64            `gorm:"type:bigint" json:"deleted_at"`
	UserID           string           `gorm:"size:20;index:idx_user_id_version;index:idx_user_id_conversation_type_conversation_id,unique" json:"-"`
	ConversationType ConversationType `gorm:"type:tinyint;index:idx_user_id_conversation_type_conversation_id,unique;index:idx_conversation_type_conversation_id;" json:"type"`
	ConversationID   string           `gorm:"size:40;index:idx_user_id_conversation_type_conversation_id,unique;index:idx_conversation_type_conversation_id;" json:"id"`
	AckSeq           int64            `gorm:"type:bigint" json:"ack_seq"`
	CleanSeq         int64            `gorm:"type:bigint" json:"clean_seq"`
	StartSeq         int64            `gorm:"type:bigint" json:"-"`
	Version          int64            `gorm:"type:bigint;index:idx_user_id_version" json:"version"`
}
