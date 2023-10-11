package model

import (
	"im/pkg/pagination"
)

type ConversationInfo struct {
	Conversation
	ConversationName    string       `json:"name"`
	ConversationFaceUrl string       `json:"face_url"`
	ReadSeq             int64        `json:"read_seq"`
	MaxSeq              int64        `json:"max_seq"`
	UnreadCount         int64        `json:"unread_count"`
	Message             *MessageInfo `json:"message"`
}

type MessageInfo struct {
	Message
	SendNickname string `json:"send_nickname"`
	SendFaceUrl  string `json:"send_face_url"`
	RecvID       string `json:"recv_id"`
	RecvName     string `json:"recv_name"`
}

type RTCInfo struct {
	ConversationType    ConversationType `json:"conversation_type"`
	SendID              string           `json:"send_id"`
	SendNickname        string           `json:"send_nickname"`
	SendFaceURL         string           `json:"send_face_url"`
	SendDeviceID        string           `json:"send_device_id"`
	RecvID              string           `json:"recv_id"`
	RecvNickname        string           `json:"recv_nickname"`
	RecvFaceURL         string           `json:"recv_face_url"`
	RecvDeviceID        string           `json:"recv_device_id"`
	RTCChannel          string           `json:"rtc_channel"`
	RTCToken            string           `json:"rtc_token"`
	RTCType             RTCType          `json:"rtc_type"`
	RTCStatus           RTCStatusType    `json:"rtc_status"`
	RTCStartTime        int64            `json:"rtc_start_time"`
	RTCRequestLimitTime int64            `json:"rtc_request_limit_time"`
	RTCRetainTime       int64            `json:"rtc_retain_time"`
	RTCUpdateTime       int64            `json:"rtc_update_time"`
}

type ConversationListReq struct {
	pagination.Pagination
	OperationID string `json:"operation_id" form:"operation_id"`
	Version     int64  `json:"version" form:"version"`
}
type ConversationListResp struct {
	pagination.Pagination
	List []ConversationInfo `json:"list"`
}

type ConversationAckSeqReq struct {
	OperationID      string           `json:"operation_id" form:"operation_id"`
	ConversationType ConversationType `json:"conversation_type" form:"conversation_type"`
	ConversationID   string           `json:"conversation_id" form:"conversation_id"`
	AckSeq           int64            `json:"ack_seq" form:"ack_seq"`
}
type ConversationAckSeqResp struct {
	ConversationType ConversationType `json:"conversation_type"`
	ConversationID   string           `json:"conversation_id"`
	AckSeq           int64            `json:"ack_seq"`
}
type ConversationReadSeqResp struct {
	ConversationType ConversationType `json:"conversation_type"`
	ConversationID   string           `json:"conversation_id"`
	ReadSeq          int64            `json:"read_seq"`
}

type MessagePullReq struct {
	OperationID      string           `json:"operation_id" form:"operation_id"`
	ConversationType ConversationType `json:"conversation_type" form:"conversation_type"`
	RecvID           string           `json:"recv_id" form:"recv_id"`
	Seq              int64            `json:"seq" form:"seq"`
	SeqList          []int64          `json:"seq_list" form:"seq_list"`
	PageSize         int              `json:"page_size" form:"page_size"`
}
type MessagePullResp struct {
	List []MessageInfo `json:"list"`
}

type MessageSendReq struct {
	OperationID      string           `json:"operation_id" form:"operation_id"`
	ConversationType ConversationType `json:"conversation_type" form:"conversation_type"`
	RecvID           string           `json:"recv_id" form:"recv_id"`
	Type             MessageType      `json:"type" form:"type"`
	Content          string           `json:"content" form:"content"`
	ClientMsgID      string           `json:"client_msg_id" form:"client_msg_id"`
	AtList           []string         `json:"at_list" form:"at_list"`
}
type MessageSendResp struct {
	MessageInfo
}

type MessageMultiSendReq struct {
	OperationID      string           `json:"operation_id" form:"operation_id"`
	ConversationType ConversationType `json:"conversation_type" form:"conversation_type"`
	Type             MessageType      `json:"type" form:"type"`
	Content          string           `json:"content" form:"content"`
	RecvIDList       []string         `json:"recv_id_list"`
	MessageList      []struct {
		ClientMsgID string `json:"client_msg_id"`
		RecvID      string `json:"recv_id"`
	} `json:"message_list"`
}
type MessageMultiSendResp struct {
	List []MessageInfo `json:"list"`
}

type MessageForwardReq struct {
	OperationID string `json:"operation_id"`
	MessageList []struct {
		ConversationType ConversationType `json:"conversation_type"`
		ClientMsgID      string           `json:"client_msg_id"`
		Type             MessageType      `json:"type"`
		Content          string           `json:"content"`
		RecvID           string           `json:"recv_id"`
	} `json:"message_list"`
}
type MessageForwardResp struct {
	List []MessageInfo `json:"list"`
}

type MessageChangeReq struct {
	OperationID string            `json:"operation_id" form:"operation_id"`
	MsgIDList   []string          `json:"msg_id_list" form:"msg_id_list"`
	Status      MessageStatusType `json:"status" form:"status"`
}
type MessageChangeResp struct {
	MessageInfo
}

type MessageClearReq struct {
	OperationID      string           `json:"operation_id" form:"operation_id"`
	Type             int              `json:"type" form:"type"`
	ConversationType ConversationType `json:"conversation_type" form:"conversation_type"`
	ConversationID   string           `json:"conversation_id" form:"conversation_id"`
	MaxSeq           int64            `json:"max_seq" form:"max_seq"`
}
type MessageClearResp struct {
	Type             int              `json:"type"`
	ConversationType ConversationType `json:"conversation_type"`
	ConversationID   string           `json:"conversation_id"`
	MaxSeq           int64            `json:"max_seq"`
}

type RTCInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id"`
}
type RTCInfoResp struct {
	*RTCInfo
}

type RTCReq struct {
	OperationID      string           `json:"operation_id" form:"operation_id"`
	DeviceID         string           `json:"device_id" form:"device_id"`
	ConversationType ConversationType `json:"conversation_type"`
	RecvID           string           `json:"recv_id" form:"recv_id"`
	RTCType          RTCType          `json:"rtc_type" form:"rtc_type"`
}
type RTCResp struct {
	*RTCInfo
}

type RTCOperateReq struct {
	OperationID   string           `json:"operation_id" form:"operation_id"`
	DeviceID      string           `json:"device_id" form:"device_id"`
	OperationType RTCOperationType `json:"operation_type" form:"operation_type"`
	RTCType       RTCType          `json:"rtc_type" form:"rtc_type"`
}
type RTCOperateResp struct {
	*RTCInfo
}

type RTCUpdateReq struct {
	OperationID string `json:"operation_id" form:"operation_id"`
	DeviceID    string `json:"device_id" form:"device_id"`
}
type RTCUpdateResp struct {
}
