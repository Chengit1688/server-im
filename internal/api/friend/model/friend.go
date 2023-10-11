package model

import "im/internal/api/user/model"

const (
	NotBlack = 2
	InBlack  = 1
)

type Friend struct {
	ID             int64  `gorm:"column:id;primary_key" json:"id"`
	CreatedAt      int64  `gorm:"column:created_at"`
	UpdatedAt      int64  `gorm:"column:updated_at"`
	OwnerUserID    string `gorm:"column:owner_user_id;index;size:64" json:"owner_user_id"`
	FriendUserID   string `gorm:"column:friend_user_id;index;size:64" json:"friend_user_id"`
	FriendLabel    string `gorm:"column:friend_label;index;size:64" json:"friend_label"`
	Remark         string `gorm:"column:remark;type:char(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"remark"`
	AddSource      int32  `gorm:"column:add_source" json:"add_source"`
	OperatorUserID string `gorm:"column:operator_user_id;size:64" json:"operator_user_id"`
	Version        int64  `gorm:"column:version" json:"version"`
	BlackStatus    int    `gorm:"column:black_status;default:2" json:"black_status"`
	Status         int64  `gorm:"column:status;default:1;" json:"status"`
	Ex             string `gorm:"column:ex;size:1024" json:"ex"`
}

func (Friend) TableName() string {
	return "friends"
}

type FriendRequest struct {
	ID            int64  `gorm:"column:id;primary_key" json:"id"`
	CreatedAt     int64  `gorm:"column:created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at"`
	FromUserID    string `gorm:"column:from_user_id;index;size:64" json:"owner_user_id"`
	ToUserID      string `gorm:"column:to_user_id;index;size:64" json:"to_user_id"`
	HandleResult  int32  `gorm:"column:handle_result" json:"handle_result"`
	ReqMsg        string `gorm:"column:req_msg;size:255" json:"req_msg"`
	HandlerUserID string `gorm:"column:handler_user_id;size:64" json:"handler_user_id"`
	HandleMsg     string `gorm:"column:handle_msg;size:255" json:"handle_msg"`
	Ex            string `gorm:"column:ex;size:1024" json:"ex"`
	Status        int    `gorm:"column:status;default:0" json:"status"`
	Remark        string `gorm:"column:remark;size:50" json:"remark"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

type Black struct {
	ID             int64      `gorm:"column:id;primary_key" json:"id"`
	CreatedAt      int64      `gorm:"column:created_at"`
	UpdatedAt      int64      `gorm:"column:updated_at"`
	OwnerUserID    string     `gorm:"column:owner_user_id;primary_key;size:64" json:"owner_user_id"`
	BlockUserID    string     `gorm:"column:block_user_id;primary_key;size:64" json:"block_user_id"`
	BlockUser      model.User `gorm:"foreignKey:UserID;references:BlockUserID"`
	AddSource      int32      `gorm:"column:add_source" json:"add_source"`
	OperatorUserID string     `gorm:"column:operator_user_id;size:64" json:"operator_user_id"`
	Status         int64      `gorm:"column:status;default:1;" json:"status"`
	Ex             string     `gorm:"column:ex;size:1024" json:"ex"`
}

func (Black) TableName() string {
	return "blacks"
}

type FriendLabel struct {
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	UserId     string `gorm:"column:user_id;size:64;index:user_id,label_id" json:"user_id"`
	LabelId    string `gorm:"column:label_id;size:20;index:user_id,label_id" json:"label_id"`
	LabelName  string `gorm:"column:label_name;type:varchar(255)" json:"label_name"`
	CreateTime int64  `gorm:"column:create_time;" json:"create_time"`
}

func (FriendLabel) TableName() string {
	return "friend_label"
}
