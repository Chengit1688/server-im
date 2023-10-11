package model

import (
	"im/internal/cms_api/config/model"

	WalletModel "im/internal/cms_api/wallet/model"
	"im/pkg/pagination"
)

type GetWalletReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
}

type GetWalletResp struct {
	PayPasswdSet int                    `json:"pay_passwd_set"`
	Balance      int64                  `json:"balance"`
	Deposit      model.GetDepositeResp  `json:"deposit"`
	Withdraw     model.WithdrawApigResp `json:"withdraw"`
}

type RedpackSingleSendReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	Amount      int64  `json:"amount"  binding:"required,gt=0,lte=200" msg:"amount required"`
	RecvID      string `json:"recv_id" binding:"required"`
	PayPasswd   string `json:"pay_passwd" binding:"required"`
	SendID      string `json:"send_id"`
	Remark      string `json:"remark"`
	MsgType     int64  `json:"msg_type"`
}

type RedpackSingleSendResp struct {
	Amount          int64                           `json:"amount"`
	RedpackSingleID int64                           `json:"redpack_single_id"`
	Status          WalletModel.StatusRedpackSingle `json:"status"`
}

type RedpackSingleRecvReq struct {
	OperationID     string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	RedpackSingleID int64  `json:"redpack_single_id" binding:"required,gte=1"`
}

type RedpackSingleRecvResp struct {
	Amount          int64                           `json:"amount"`
	RedpackSingleID int64                           `json:"redpack_single_id"`
	Status          WalletModel.StatusRedpackSingle `json:"status"`
	RecvAt          *int64                          `json:"recv_at"`
	Remark          string                          `json:"remark"`
}

type RedpackSingleGetInfoReq struct {
	OperationID     string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	RedpackSingleID int64  `json:"redpack_single_id" form:"redpack_single_id" binding:"required,gte=1"`
}

type RedpackSingleGetInfoResp struct {
	Amount          int64                           `json:"amount"`
	RedpackSingleID int64                           `json:"redpack_single_id"`
	Status          WalletModel.StatusRedpackSingle `json:"status"`
	RecvAt          *int64                          `json:"recv_at"`
}

type WalletSetPayPassReqs struct {
	OperationID     string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	PayPasswd       string `json:"pay_passwd" form:"pay_passwd" binding:"required,eqfield=ConfirmPassword,min=6,max=6,number"  msg:"请正确输入支付密码,6位数字"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=PayPasswd,min=6,max=6,number"  msg:"请再次确认支付密码是否正确"`
}

type RedpackSingleMessagePush struct {
	ConversationID   string `json:"conversation_id"`
	Timestamp        int64  `json:"timestamp"`
	SenderNickname   string `json:"sender_nickname"`
	ReceiverNickname string `json:"receiver_nickname"`
	RedpackID        int64  `json:"redpack_id"`
	Type             int64  `json:"type"`
}

type RedpackGroupMessagePush struct {
	ConversationID   string `json:"conversation_id"`
	Timestamp        int64  `json:"timestamp"`
	SenderNickname   string `json:"sender_nickname"`
	ReceiverNickname string `json:"receiver_nickname"`
	GroupId          string `json:"group_id"`
	RedpackGroupID   int64  `json:"redpack_group_id"`
	Type             int64  `json:"type"`
}

type WithdrawCommitReq struct {
	OperationID string                         `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Columns     model.WithdrawConfigColumnList `json:"columns"`
}

type WithdrawCommitResp struct {
	ID        int64 `json:"id"`
	BillingID int64 `json:"billing_id"`
}

type BillingRecordsListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id"`
	Direction   string `json:"direction" form:"direction" binding:"omitempty,oneof=all in out"  msg:"只支持int/out"`
	pagination.Pagination
}

type BillingRecordsListResp struct {
	List []BillingRecordsListItem `json:"list"`
	pagination.Pagination
}

type BillingRecordsListItem struct {
	Amount    int64                         `json:"amount"`
	Name      string                        `json:"name"`
	Type      WalletModel.TypeBillingRecord `json:"type"`
	CreatedAt int64                         `json:"created_at"`
	UpdatedAt int64                         `json:"updated_at"`
}

type RedpackGroupSendReq struct {
	OperationID string                       `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	Amount      int64                        `json:"amount"  binding:"required,gt=1" msg:"amount required"`
	PreAmount   int64                        `json:"pre_amount"  binding:"gt=1" msg:"普通红包的每个红包的金额"`
	PayPasswd   string                       `json:"pay_passwd" binding:"required"`
	SendID      string                       `json:"send_id"`
	GroupID     string                       `json:"group_id" binding:"required,min=1"`
	Count       int64                        `json:"count"  binding:"required,number"`
	MsgType     int64                        `json:"msg_type"  binding:"required,number"`
	Type        WalletModel.TypeRedpackGroup `json:"type"  binding:"required,number,oneof=1 2"`
	Remark      string                       `json:"remark"`
}

type RedpackGroupSendResp struct {
	Amount         int64                           `json:"amount"`
	RedpackGroupID int64                           `json:"redpack_group_id"`
	Status         WalletModel.StatusRedpackSingle `json:"status"`
}

type RedpackGroupRecvReq struct {
	OperationID    string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	RedpackGroupID int64  `json:"redpack_group_id" binding:"required,gte=1"`
	GroupID        string `json:"group_id" binding:"required,gte=1"`
}

type RedpackGroupRecvResp struct {
	Amount         int64                           `json:"amount"`
	RedpackGroupID int64                           `json:"redpack_group_id"`
	Status         WalletModel.StatusRedpackSingle `json:"status"`
	SendAt         int64                           `json:"send_at"`
	RecvAt         int64                           `json:"recv_at"`
	Type           WalletModel.TypeRedpackGroup    `json:"type"`
	Remark         string                          `json:"remark"`
}

type RedpackGroupGetInfoReq struct {
	OperationID    string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"operation_id required"`
	RedpackGroupID int64  `json:"redpack_group_id" form:"redpack_group_id" binding:"required,gte=1"`
	GroupID        string `json:"group_id" form:"group_id" binding:"required,gte=1"`
}

type RedpackGroupGetInfoResp struct {
	Amount          int64                           `json:"amount"`
	RedpackSingleID int64                           `json:"redpack_single_id"`
	Status          WalletModel.StatusRedpackSingle `json:"status"`
	RecvAt          *int64                          `json:"recv_at"`
}
