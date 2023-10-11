package model

type MessageType int

const (
	MessageNone  MessageType = 0
	MessageText  MessageType = 1
	MessageFace  MessageType = 2
	MessageImage MessageType = 3
	MessageVoice MessageType = 4
	MessageVideo MessageType = 5
	MessageFile  MessageType = 6
	MessageCard  MessageType = 7
	MessageAt    MessageType = 8
	MessageQuote MessageType = 9
	MessageLink  MessageType = 10

	MessageOperation MessageType = 100
	MessageRead      MessageType = 101
	MessageRevoke    MessageType = 102
	MessageDelete    MessageType = 103

	MessageFriend          MessageType = 200
	MessageFriendAddNotify MessageType = 201

	MessageGroup                   MessageType = 300
	MessageGroupCreateNotify       MessageType = 301
	MessageGroupAddMemberNotify    MessageType = 302
	MessageGroupDeleteNotify       MessageType = 303
	MessageGroupSetAdminNotify     MessageType = 304
	MessageGroupUnsetAdminNotify   MessageType = 305
	MessageGroupOneMuteNotify      MessageType = 306
	MessageGroupOneUnmuteNotify    MessageType = 307
	MessageGroupAllMuteNotify      MessageType = 308
	MessageGroupAllUnmuteNotify    MessageType = 309
	MessageGroupTransferNotify     MessageType = 310
	MessageGroupNotifyChangeNotify MessageType = 311

	MessageRTCNotify MessageType = 400
)

func MessageTypeString(mt MessageType) string {
	switch mt {
	case MessageText:
		return "文本"
	case MessageFace:
		return "表情"
	case MessageImage:
		return "图片"
	case MessageVoice:
		return "语音"
	case MessageVideo:
		return "视频"
	case MessageFile:
		return "文件"
	case MessageCard:
		return "名片"
	case MessageAt:
		return "at 消息"
	case MessageQuote:
		return "引用消息"
	}
	return "未知"
}

type MessageStatusType int

const (
	MessageStatusTypeNotRead MessageStatusType = 0
	MessageStatusTypeRead    MessageStatusType = 1
	MessageStatusTypeRevoke  MessageStatusType = 2
	MessageStatusTypeDelete  MessageStatusType = 3
)

type MessageContentTimeType int

const (
	MessageContentTimeTypeOneHour MessageContentTimeType = 1
	MessageContentTimeTypeOneDay  MessageContentTimeType = 2
	MessageContentTimeTypeForever MessageContentTimeType = 3
)

type Message struct {
	ID               uint              `gorm:"primaryKey" json:"-"`
	MsgID            string            `gorm:"size:64;index:idx_msg_id,unique" json:"msg_id"`
	ClientMsgID      string            `gorm:"size:64;" json:"client_msg_id"`
	ConversationType ConversationType  `gorm:"type:tinyint;index:idx_conversation_type_conversation_id_seq" json:"conversation_type"`
	ConversationID   string            `gorm:"size:40;index:idx_conversation_type_conversation_id_seq" json:"conversation_id"`
	SendID           string            `gorm:"size:20" json:"send_id"`
	SendTime         int64             `gorm:"type:bigint;index:idx_send_time" json:"send_time"`
	Type             MessageType       `gorm:"type:smallint" json:"type"`
	Content          string            `gorm:"type:text" json:"content"`
	Status           MessageStatusType `gorm:"type:tinyint" json:"status"`
	Seq              int64             `gorm:"type:bigint;index:idx_conversation_type_conversation_id_seq,unique" json:"seq"`
	Role             string            `gorm:"column:role;type:varchar(255);default:'';" json:"role"`
}

type MessageBeOperator struct {
	BeOperatorID       string `json:"be_operator_id"`
	BeOperatorNickname string `json:"be_operator_nickname"`
}

type MessageContent struct {
	OperatorID       string                 `json:"operator_id"`
	OperatorNickname string                 `json:"operator_nickname"`
	OperatorFaceUrl  string                 `json:"operator_face_url"`
	BeOperatorList   []MessageBeOperator    `json:"be_operator_list"`
	TimeType         MessageContentTimeType `json:"time_type"`
	MsgIDList        []string               `json:"msg_id_list"`
	Content          string                 `json:"content"`
}

type MsgContent struct {
	Text string `json:"text"`
}

type MessageContentClient struct {
	Text      string `json:"text"`
	ImageInfo struct {
		ImageURL string `json:"image_url"`
	} `json:"image_info"`
	AudioInfo struct {
		FileURL string `json:"file_url"`
	}
	VideoInfo struct {
		FileURL string `json:"file_url"`
	}
	FileInfo struct {
		FileURL string `json:"file_url"`
	}
	CardInfo struct {
		Nickname string `json:"nick_name"`
	}
}
