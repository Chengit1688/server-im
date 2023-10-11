package model

import (
	apiChatModel "im/internal/api/chat/model"
	"im/pkg/pagination"
)

type MessageHistoryListReq struct {
	OperationID      string                        `json:"operation_id" form:"operation_id" binding:"required"`
	ConversationType apiChatModel.ConversationType `json:"conversation_type" form:"conversation_type"`
	SendID           string                        `json:"send_id" form:"send_id"`
	RecvID           string                        `json:"recv_id" form:"recv_id"`
	Type             apiChatModel.MessageType      `json:"type" form:"type"`
	Content          string                        `json:"content" form:"content"`
	StartTime        int64                         `json:"start_time" form:"start_time"`
	EndTime          int64                         `json:"end_time" form:"end_time"`
	Export           bool                          `json:"export" form:"export"`
	pagination.Pagination
}

type MessageHistoryListResp struct {
	List []apiChatModel.MessageInfo `json:"list"`
	pagination.Pagination
}

type MessageClearReq struct {
	OperationID string    `json:"operation_id" form:"operation_id"`
	Type        ClearType `json:"type" form:"type"`
	TargetID    string    `json:"target_id" form:"target_id"`
}
type MessageClearResp struct {
	Type     ClearType `json:"type"`
	TargetID string    `json:"target_id"`
}

type GetMultiSendPagingReq struct {
	OperationID      string `json:"operation_id" form:"operation_id" binding:"required"`
	Content          string `json:"content" form:"content"`
	SenderID         string `json:"sender_id" form:"sender_id"`
	SenderNickname   string `json:"sender_nickname" form:"sender_nickname"`
	Operate          string `json:"operate" form:"operate"`
	OperateTimeStart int64  `json:"operate_time_start" form:"operate_time_start"`
	OperateTimeEnd   int64  `json:"operate_time_end" form:"operate_time_end"`
	pagination.Pagination
}

type GetMultiSendPagingResp struct {
	List []GetMultiSendPagingItem `json:"list"`
	pagination.Pagination
}

type GetMultiSendPagingItem struct {
	ID              int    `json:"id"`
	SenderIDs       string `json:"sender_ids"`
	SenderNicknames string `json:"sender_nicknames"`
	Content         string `json:"content"`
	Operate         string `json:"operate"`
	CreatedAt       int64  `json:"created_at"`
}

type MultiSendReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required"`
	Content     string `json:"content" binding:"required" msg:"消息内容不能为空"`
	SenderIDs   string `json:"sender_ids" binding:"required" msg:"发送者ID不能为空"`
}
