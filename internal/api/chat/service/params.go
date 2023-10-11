package service

type MessageMultiSendReq struct {
	OperationID string `json:"operation_id" form:"operation_id"`

	ConversationType int      `json:"conversation_type" form:"conversation_type"`
	Type             int      `json:"type" form:"type"`
	Content          string   `json:"content" form:"content"`
	RecvIDList       []string `json:"recv_id_list"`
	MessageList      []struct {
		ClientMsgID string `json:"client_msg_id"`
		RecvID      string `json:"recv_id"`
	} `json:"message_list"`
}

type MessageMultiSendResp struct {
	List []MessageInfo `json:"list"`
}

type MessageInfo struct {
	Message
	SendNickname string `json:"send_nickname"`
	SendFaceUrl  string `json:"send_face_url"`
	RecvID       string `json:"recv_id"`
	RecvName     string `json:"recv_name"`
}

type Message struct {
	ID               uint   `gorm:"primaryKey" json:"-"`
	MsgID            string `gorm:"size:64;index:idx_msg_id,unique" json:"msg_id"`
	ClientMsgID      string `gorm:"size:64;" json:"client_msg_id"`
	ConversationType int    `gorm:"type:tinyint;index:idx_conversation_type_conversation_id_seq" json:"conversation_type"`
	ConversationID   string `gorm:"size:40;index:idx_conversation_type_conversation_id_seq" json:"conversation_id"`
	SendID           string `gorm:"size:20" json:"send_id"`
	SendTime         int64  `gorm:"type:bigint;index:idx_send_time" json:"send_time"`
	Type             int    `gorm:"type:smallint" json:"type"`
	Content          string `gorm:"type:text" json:"content"`
	Status           int    `gorm:"type:tinyint" json:"status"`
	Seq              int64  `gorm:"type:bigint;index:idx_conversation_type_conversation_id_seq,unique" json:"seq"`
}
