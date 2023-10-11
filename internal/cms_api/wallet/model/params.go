package model

import (
	configModel "im/internal/cms_api/config/model"
	"im/pkg/pagination"
)

type BillingRecordsListItemResp struct {
	ID               int64             `json:"id"`
	SenderID         string            `json:"sender_id"`
	SenderNickName   string            `json:"sender_nick_name"`
	ReceiverID       string            `json:"receiver_id"`
	ReceiverNickName string            `json:"receiver_nick_name"`
	CreatedAt        int64             `json:"created_at"`
	Type             TypeBillingRecord `json:"type"`
	Amount           int64             `json:"amount"`
	ChangeBefore     int64             `json:"change_before"`
	ChangeAfter      int64             `json:"change_after"`
	Note             string            `json:"note"`
}

type BillingRecordsListResp struct {
	List     []BillingRecordsListItemResp `json:"list"`
	Count    int64                        `json:"count"`
	Page     int                          `json:"page"`
	PageSize int                          `json:"page_size"`
}

type BillingRecordsListReq struct {
	OperationID      string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID               int64  `json:"id" form:"id"`
	SenderID         string `json:"sender_id" form:"sender_id"`
	ReceiverID       string `json:"receiver_id" form:"receiver_id"`
	Type             int    `json:"type" form:"type"`
	CreatedTimeStart int64  `json:"created_time_start" form:"created_time_start"`
	CreatedTimeEnd   int64  `json:"created_time_end" form:"created_time_end"`
	Direction        string `json:"direction" form:"direction" binding:"omitempty,oneof=in out"  msg:"只支持int/out"`
	pagination.Pagination
}

type WalletChangeAmountReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id" form:"user_id" binding:"required"`
	Amount      int64  `json:"amount" form:"amount" binding:"required"`
	Note        string `json:"note" form:"note" binding:"required"`
}

type RedpackSingleRecordsListItemResp struct {
	ID               int64               `json:"id"`
	SenderID         string              `json:"sender_id"`
	SenderNickName   string              `json:"sender_nick_name"`
	ReceiverID       string              `json:"receiver_id"`
	ReceiverNickName string              `json:"receiver_nick_name"`
	SendAt           int64               `json:"send_at"`
	RecvAt           *int64              `json:"recv_at"`
	Status           StatusRedpackSingle `json:"status"`
	Amount           int64               `json:"amount"`
}

type RedpackSingleRecordsListResp struct {
	List     []RedpackSingleRecordsListItemResp `json:"list"`
	Count    int64                              `json:"count"`
	Page     int                                `json:"page"`
	PageSize int                                `json:"page_size"`
}

type RedpackSingleRecordsListReq struct {
	OperationID   string              `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	SenderID      string              `json:"sender_id" form:"sender_id"`
	ReceiverID    string              `json:"receiver_id" form:"receiver_id"`
	Status        StatusRedpackSingle `json:"status" form:"status"`
	SendTimeStart int64               `json:"send_time_start" form:"send_time_start"`
	SendTimeEnd   int64               `json:"send_time_end" form:"send_time_end"`
	RecvTimeStart int64               `json:"recv_time_start" form:"recv_time_start"`
	RecvTimeEnd   int64               `json:"recv_time_end" form:"recv_time_end"`
	pagination.Pagination
}

type RedpackGroupRecordsListItemResp struct {
	Id             int    `json:"id"`
	GroupName      string `json:"group_name"`
	SenderNickName string `json:"sender_nick_name"`
	SenderUserId   string `json:"sender_user_id"`
	Type           int64  `json:"type"`
	SendAt         int64  `json:"send_at"`
	Status         int    `json:"status"`
	Total          int64  `json:"total"`
	Amount         int64  `json:"amount"`
}

type RedpackGroupRecordsListResp struct {
	List []RedpackGroupRecordsListItemResp `json:"list"`
	pagination.Pagination
}

type RedpackGroupRecordsListReq struct {
	OperationID   string              `json:"operation_id" binding:"required" msg:"操作ID不能为空"`
	GroupName     string              `json:"group_name"  binding:"omitempty,min=1" msg:"请输入群名称"`
	SenderID      string              `json:"sender_id" binding:"omitempty,min=1" msg:"请输入userId"`
	Status        StatusRedpackSingle `json:"status" binding:"omitempty,oneof=1 2 3" msg:"请选择红包状态"`
	SendTimeStart int                 `json:"send_time_start" binding:"omitempty,number" msg:"请选择红包发送时间"`
	SendTimeEnd   int                 `json:"send_time_end" binding:"omitempty,number" msg:"请选择红包发送的结束时间"`
	pagination.Pagination
}
type WithdrawRecordsListReq struct {
	OperationID      string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	BillingID        int64  `json:"billing_id" form:"billing_id"`
	UserID           string `json:"user_id" form:"user_id"`
	NickName         string `json:"nick_name" form:"nick_name"`
	Status           int    `json:"status" form:"status"`
	CreatedTimeStart int64  `json:"created_time_start" form:"created_time_start"`
	CreatedTimeEnd   int64  `json:"created_time_end" form:"created_time_end"`
	IsDone           string `json:"is_done" form:"is_done" binding:"omitempty,oneof=yes no"  msg:"yes/no"`
	pagination.Pagination
}

type WithdrawRecordsListItemResp struct {
	ID        int64  `json:"id"`
	BillingID int64  `json:"billing_id" form:"billing_id"`
	UserID    string `json:"user_id" form:"user_id"`
	NickName  string `json:"nick_name" form:"nick_name"`
	Status    int    `json:"status" form:"status"`
	IsDone    string `json:"is_done" form:"is_done" binding:"omitempty,oneof=yes no"  msg:"yes/no"`
	CreatedAt int64  `json:"created_at"`
	Amount    int64  `json:"amount"`
}

type WithdrawRecordsListResp struct {
	List     []WithdrawRecordsListItemResp `json:"list"`
	Count    int64                         `json:"count"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"page_size"`
}

type GetWithdrawRecordsNotDoneReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetWithdrawRecordsDescribeResp struct {
	ID      int64                                `json:"id"`
	Columns configModel.WithdrawConfigColumnList `json:"columns"`
	Note    string                               `json:"note"`
}

type SetWithdrawRecordsStatusReq struct {
	OperationID string         `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	ID          int64          `json:"id" binding:"required"`
	Status      StatusWithdraw `json:"status" binding:"required,oneof=1 2"  msg:"1/2"`
	Note        string         `json:"note"`
}

type SetWithdrawRecordsStatusResp struct {
	ID     int64          `json:"id"`
	Status StatusWithdraw `json:"status"`
	Note   string         `json:"note"`
}

type WalletSetPayPassReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	PayPasswd   string `json:"pay_passwd" form:"pay_passwd" binding:"required,min=6,max=6,number"  msg:"请正确输入支付密码,6位数字"`
	UserID      string `json:"user_id" binding:"required,min=10,max=15"`
}

type RedpackGroupReturnsRecords struct {
	GroupRecordsID  int64  `json:"group_records_id"`
	SendAt          int64  `json:"send_at"`
	SenderID        string `json:"sender_id"`
	GroupID         string `json:"group_id"`
	Amount          int64  `json:"amount"`
	RecvAmount      int64  `json:"recv_amount"`
	RemainderAmount int64  `json:"remainder_amount"`
}
