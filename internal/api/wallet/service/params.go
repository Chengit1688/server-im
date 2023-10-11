package service

import (
	"im/pkg/pagination"
)

type GetWalletReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
}

type RedpackSingleSendResp struct {
	Amount          int64 `json:"amount"`
	RedpackSingleID int64 `json:"redpack_single_id"`
	Status          int   `json:"status"`
}

type RedpackSingleRecvResp struct {
	Amount          int64  `json:"amount"`
	RedpackSingleID int64  `json:"redpack_single_id"`
	Status          int    `json:"status"`
	RecvAt          *int64 `json:"recv_at"`
	Remark          string `json:"remark"`
}

type BillingRecordsListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id"`
	Direction   string `json:"direction" form:"direction" binding:"omitempty,oneof=in out"  msg:"只支持int/out"`
	pagination.Pagination
}

type BillingRecordsListResp struct {
	List []BillingRecordsListItem `json:"list"`
	pagination.Pagination
}

type BillingRecordsListItem struct {
	Amount    int64  `json:"amount"`
	Name      string `json:"name"`
	Type      int    `json:"type"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type RedpackGroupSendReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	Amount      int64  `json:"amount"  binding:"required,gt=1" msg:"amount required"`
	PreAmount   int64  `json:"pre_amount"  binding:"gt=1" msg:"普通红包的每个红包的金额"`
	PayPasswd   string `json:"pay_passwd" binding:"required"`
	SendID      string `json:"send_id"`
	GroupID     string `json:"group_id" binding:"required,min=1"`
	Count       int64  `json:"count"  binding:"required,number"`
	MsgType     int64  `json:"msg_type"  binding:"required,number"`
	Type        int    `json:"type"  binding:"required,number,oneof=1 2"`
	Remark      string `json:"remark"`
}

type RedpackGroupSendResp struct {
	Amount         int64 `json:"amount"`
	RedpackGroupID int64 `json:"redpack_group_id"`
	Status         int   `json:"status"`
}

type RedpackGroupRecvResp struct {
	Amount         int64  `json:"amount"`
	RedpackGroupID int64  `json:"redpack_group_id"`
	Status         int    `json:"status"`
	SendAt         int64  `json:"send_at"`
	RecvAt         int64  `json:"recv_at"`
	Type           int    `json:"type"`
	Remark         string `json:"remark"`
}
