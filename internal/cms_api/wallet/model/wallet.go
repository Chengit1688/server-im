package model

import (
	apiGroupModel "im/internal/api/group/model"
	apiUserModel "im/internal/api/user/model"
	cmsUserModel "im/internal/cms_api/admin/model"
	"im/pkg/db"
)

type StatusRedpackSingle int

const StatusRedpackSingleRecv StatusRedpackSingle = 1
const StatusRedpackSingleRecvd StatusRedpackSingle = 2
const StatusRedpackSingleReturn StatusRedpackSingle = 3

type StatusWithdraw int

const StatusWithdrawPermit StatusWithdraw = 1
const StatusWithdrawRefused StatusWithdraw = 2

type TypeRedpackGroup int

const TypeRedpackGroupRandom TypeRedpackGroup = 1
const TypeRedpackGroupNormal TypeRedpackGroup = 2

type TypeBillingRecord int

var (
	TypeBillingRecordCmsDeposit          TypeBillingRecord = 1
	TypeBillingRecordCmsChange           TypeBillingRecord = 2
	TypeBillingRecordSendRedPackSingle   TypeBillingRecord = 3
	TypeBillingRecordRecvRedPackSingle   TypeBillingRecord = 4
	TypeBillingRecordSendRedPackGroup    TypeBillingRecord = 5
	TypeBillingRecordRecvRedPackGroup    TypeBillingRecord = 6
	TypeBillingRecordRedPackGroupReturn  TypeBillingRecord = 7
	TypeBillingRecordRedPackSingleReturn TypeBillingRecord = 8
	TypeBillingRecordWithdraw            TypeBillingRecord = 9
	TypeBillingRecordWithdrawSuccess     TypeBillingRecord = 10
	TypeBillingRecordWithdrawFailed      TypeBillingRecord = 11
	TypeBillingRecordWithdrawRollback    TypeBillingRecord = 12
	TypeBillingRecordWithdrawSignAward   TypeBillingRecord = 13
)

type BillingRecords struct {
	db.CommonModel
	SenderID       string              `gorm:"column:sender_id;index:idx_sender_id;size:32" json:"sender_id"`
	SenderApiUser  apiUserModel.User   `gorm:"foreignKey:UserID;references:SenderID"`
	SenderCmsUser  cmsUserModel.Admin  `gorm:"foreignKey:UserID;references:SenderID"`
	ReceiverID     string              `gorm:"column:receiver_id;index:idx_receiver_id;size:32" json:"receiver_id"`
	Receiver       apiUserModel.User   `gorm:"foreignKey:UserID;references:ReceiverID"`
	Type           TypeBillingRecord   `gorm:"column:type;" json:"type"`
	Amount         int64               `gorm:"column:amount;" json:"amount"`
	ChangeBefore   int64               `gorm:"column:change_before;"  json:"change_before"`
	ChangeAfter    int64               `gorm:"column:change_after;"  json:"change_after"`
	Note           string              `gorm:"column:note;" json:"note"`
	GroupID        string              `gorm:"column:group_id;" json:"group_id"`
	RedpackGroupID int64               `gorm:"column:redpack_group_id;" json:"redpack_group_id"`
	Group          apiGroupModel.Group `gorm:"foreignKey:GroupId;references:GroupID"`
}

func (d *BillingRecords) TableName() string {
	return "billing_records"
}

type RedpackSingleRecords struct {
	ID         int64               `gorm:"column:id;primarykey"`
	SendAt     int64               `gorm:"column:send_at;index:idx_send_at"`
	RecvAt     *int64              `gorm:"column:recv_at;index:idx_recv_at"`
	SenderID   string              `gorm:"column:sender_id;index:idx_sender_id;size:32" json:"sender_id"`
	Remark     string              `gorm:"column:remark;size:120;default:恭喜发财，大吉大利" json:"remark"`
	Sender     apiUserModel.User   `gorm:"foreignKey:UserID;references:SenderID"`
	ReceiverID string              `gorm:"column:receiver_id;index:idx_receiver_id;size:32" json:"receiver_id"`
	Receiver   apiUserModel.User   `gorm:"foreignKey:UserID;references:ReceiverID"`
	MsgType    int64               `gorm:"column:msg_type;default:8;" json:"msg_type"`
	Status     StatusRedpackSingle `gorm:"column:status;" json:"status"`
	Amount     int64               `gorm:"column:amount;" json:"amount"`
}

func (d *RedpackSingleRecords) TableName() string {
	return "redpack_single_records"
}

type WithdrawRecords struct {
	db.CommonModel
	BillingID     int64             `gorm:"column:billing_id" json:"billing_id"`
	BillingRecord BillingRecords    `gorm:"foreignKey:ID;references:BillingID"`
	UserID        string            `gorm:"column:user_id;index:idx_user_id;size:32" json:"user_id"`
	User          apiUserModel.User `gorm:"foreignKey:UserID;references:UserID"`
	Amount        int64             `gorm:"column:amount;" json:"amount"`
	Status        *StatusWithdraw   `gorm:"column:status;" json:"status"`
	Columns       string            `gorm:"column:columns;" json:"columns"`
	Note          string            `gorm:"column:note;" json:"note"`
}

func (d *WithdrawRecords) TableName() string {
	return "withdraw_records"
}

type RedpackGroupRecords struct {
	ID              int64               `gorm:"column:id;primarykey"`
	SendAt          int64               `gorm:"column:send_at;index:idx_send_at"`
	SenderID        string              `gorm:"column:sender_id;index:idx_sender_id;size:32" json:"sender_id"`
	Sender          apiUserModel.User   `gorm:"foreignKey:UserID;references:SenderID"`
	GroupID         string              `gorm:"column:group_id;index:idx_group_id;size:20"`
	Group           apiGroupModel.Group `gorm:"foreignKey:GroupId;references:GroupID"`
	Amount          int64               `gorm:"column:amount;default:0;" json:"amount"`
	RemainderAmount int64               `gorm:"column:remainder_amount;default:0;" json:"remainder_amount"`
	ReceiveCount    int64               `gorm:"column:receive_count;default:0;" json:"receive_count"`
	Count           int64               `gorm:"column:count;default:0;" json:"count"`
	MsgType         int64               `gorm:"column:msg_type;default:8;" json:"msg_type"`
	Remark          string              `gorm:"column:remark;size:120;default:恭喜发财，大吉大利" json:"remark"`
	Type            TypeRedpackGroup    `gorm:"column:type;default:2;" json:"type"`
	Status          StatusRedpackSingle `gorm:"column:status;default:1;" json:"status"`
}

func (d *RedpackGroupRecords) TableName() string {
	return "redpack_group_records"
}

type RedpackGroupRecvs struct {
	ID             int64               `gorm:"column:id;primarykey"`
	RecvAt         int64               `gorm:"column:recv_at;index:idx_recv_at"`
	SendAt         int64               `gorm:"column:send_at;index:idx_send_at"`
	SenderID       string              `gorm:"column:sender_id;index:idx_sender_id;size:32;" json:"sender_id"`
	UserID         string              `gorm:"column:user_id;index:idx_user_id;size:32" json:"user_id"`
	User           apiUserModel.User   `gorm:"foreignKey:UserID;references:UserID"`
	RedpackGroupID int64               `gorm:"column:redpack_group_id;index:idx_redpack_group_id" json:"redpack_group_id"`
	RedpackGroup   RedpackGroupRecords `gorm:"foreignKey:ID;references:RedpackGroupID"`
	Amount         int64               `gorm:"column:amount;" json:"amount"`
	Status         StatusRedpackSingle `gorm:"column:status;" json:"status"`
}

func (d *RedpackGroupRecvs) TableName() string {
	return "redpack_group_recvs"
}
