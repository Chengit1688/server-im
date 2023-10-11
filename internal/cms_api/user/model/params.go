package model

import (
	"im/pkg/common"
	"im/pkg/pagination"
)

type UserListItemResp struct {
	ID              int64  `json:"id"`
	UserID          string `json:"user_id"`
	Account         string `json:"account"`
	NickName        string `json:"nick_name"`
	PhoneNumber     string `json:"phone_number"`
	CreatedAt       int64  `json:"created_at"`
	LoginIp         string `json:"login_ip"`
	IPInfo          string `json:"ip_info"`
	Balance         int64  `json:"balance"`
	InviteCode      string `json:"invite_code"`
	Status          int64  `json:"status"`
	Gender          int64  `json:"gender"`
	FaceURL         string `json:"face_url"`
	LatestLoginTime int64  `json:"latest_login_time"`
	RealName        string `json:"real_name"`
	IDNo            string `json:"id_no"`
	IDFrontImg      string `json:"id_front_img"`
	IDBackImg       string `json:"id_back_img"`
	RealAuthMsg     string ` json:"real_auth_msg"`
	IsRealAuth      int    `json:"is_real_auth"`
}

type RealNameListItemResp struct {
	ID          int64  `json:"id"`
	UserID      string `json:"user_id"`
	Account     string `json:"account"`
	NickName    string `json:"nick_name"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   int64  `json:"created_at"`
	FaceURL     string `json:"face_url"`
	RealName    string `json:"real_name"`
	IDNo        string `json:"id_no"`
	IDFrontImg  string `json:"id_front_img"`
	IDBackImg   string `json:"id_back_img"`
	RealAuthMsg string ` json:"real_auth_msg"`
	IsRealAuth  int    `json:"is_real_auth"`
}

type UserListResp struct {
	List     []UserListItemResp `json:"list"`
	Count    int64              `json:"count"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type RealNameListResp struct {
	List     []RealNameListItemResp `json:"list"`
	Count    int64                  `json:"count"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

type UserListReq struct {
	OperationID       string  `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID            string  `json:"user_id" form:"user_id"`
	NickName          string  `json:"nick_name" form:"nick_name"`
	Account           string  `json:"account" form:"account"`
	PhoneNumber       string  `json:"phone_number" form:"phone_number"`
	RegisterTimeStart int64   `json:"register_time_start" form:"register_time_start"`
	RegisterTimeEnd   int64   `json:"register_time_end" form:"register_time_end"`
	Status            int64   `json:"status" form:"status"`
	LoginIp           string  `json:"login_ip" form:"login_ip"`
	Gender            int64   `json:"gender" form:"gender"`
	InviteCode        *string `json:"invite_code" form:"invite_code"`
	IsPrivilege       int64   `json:"is_privilege" form:"is_privilege"`
	IsCustomer        int64   `json:"is_customer" form:"is_customer"`
	pagination.Pagination
}

type RealNameListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id" form:"user_id"`
	IsRealAuth  int    `json:"is_real_auth"`
	pagination.Pagination
}

type UserForGroupListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	SearchKey   string `json:"search_key" form:"search_key"`
	Status      int64  `json:"status" form:"status"`
	pagination.Pagination
}

type UserAddBatchReq struct {
	OperationID string             `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Users       []UserAddBatchItem `json:"users" binding:"dive" msg:"用户信息不能为空"`
}

type UserAddBatchItem struct {
	NickName string `json:"nick_name" binding:"required,min=1,max=16" msg:"用户昵称不能为空"`
	Account  string `json:"account" binding:"required,min=4,max=16,alphanum" msg:"用户账号不能为空"`
	Password string `json:"password" binding:"required,min=6,max=16" msg:"用户密码不能为空"`
}

type UserAddBatchResp struct {
	Users []UserAddBatchItem `json:"users"`
}

type GetUserDetailsByUserIDReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id" form:"user_id" binding:"required,min=6,max=16,alphanum" msg:"用户账号不能为空"`
}

type GetUserDetailsByUserIDResp struct {
	UserID       string `json:"user_id"`
	Account      string `json:"account"`
	PhoneNumber  string `json:"phone_number"`
	FaceURL      string `json:"face_url"`
	Gender       int64  `json:"gender"`
	NickName     string `json:"nick_name"`
	Signatures   string `json:"signatures"`
	Age          int64  `json:"age"`
	RegisterTime int64  `json:"register_time"`
	Balance      int64  `json:"balance"`
	Online       int    `json:"online"`
	LoginIp      string `json:"login_ip"`
	IsPrivilege  int64  `json:"is_privilege"`
}

type UpdateUserInfoReq struct {
	OperationID string  `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string  `json:"user_id" binding:"required,min=6,max=16,number" msg:"用户ID不能为空"`
	NickName    *string `json:"nick_name" binding:"required" msg:"昵称不能省略"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,max=20" msg:"手机号限长20"`
	Age         *int64  `json:"age"`
	Signatures  *string `json:"signatures" binding:"omitempty,max=255" msg:"签名限长255"`
}

type FreezeUserReq struct {
	OperationID string   `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string   `json:"user_id"`
	UserIDs     []string `json:"user_ids"`
}

type UnFreezeUserReq struct {
	OperationID string   `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string   `json:"user_id"`
	UserIDs     []string `json:"user_ids"`
}

type FreezePushMessage struct {
	Msg string `json:"msg"`
}

type SetUserPasswordReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id" binding:"required,min=6,max=16,number" msg:"用户ID不能为空"`
	Password    string `json:"password" binding:"required,min=6,max=16" msg:"用户密码不能为空"`
}

type DMUserListItemResp struct {
	ID          int64  `json:"id"`
	UserID      string `json:"user_id"`
	Account     string `json:"account"`
	NickName    string `json:"nick_name"`
	PhoneNumber string `json:"phone_number"`
	CountDevice int64  `json:"count_device"`
	CountIP     int64  `json:"count_ip"`
}

type DMUserListResp struct {
	List     []DMUserListItemResp `json:"list"`
	Count    int64                `json:"count"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

type DMUserListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Search      string `json:"search" form:"search"`
	Status      int    `json:"status" form:"status"`
	pagination.Pagination
}

type DMDeviceListItemResp struct {
	ID        int64  `json:"id"`
	DeviceID  string `json:"device_id"`
	Platform  int64  `json:"platform"`
	CountIP   int64  `json:"count_ip"`
	CountUser int64  `json:"count_user"`
}

type DMDeviceListResp struct {
	List     []DMDeviceListItemResp `json:"list"`
	Count    int64                  `json:"count"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

type DMDeviceListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Search      string `json:"search" form:"search"`
	Status      int    `json:"status" form:"status"`
	pagination.Pagination
}

type DMIPListItemResp struct {
	ID          int64  `json:"id"`
	IP          string `json:"ip"`
	IPInfo      string `json:"ip_info"`
	CountDevice int64  `json:"count_device"`
	CountUser   int64  `json:"count_user"`
}

type DMIPListResp struct {
	List     []DMIPListItemResp `json:"list"`
	Count    int64              `json:"count"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type DMIPListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Search      string `json:"search" form:"search"`
	Status      int    `json:"status" form:"status"`
	pagination.Pagination
}

type DeviceLockReq struct {
	OperationID string   `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	DeviceIDs   []string `json:"device_ids" binding:"required"`
}

type UserSearchReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Search      string `json:"search" form:"search" binding:"required,gte=1"`
}

type UserSearchResp struct {
	List  []UserListItemResp `json:"list"`
	Count int64              `json:"count"`
}

type SetPrivilegeUserFreezeReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	IsFreeze    int    `json:"is_freeze" binding:"required,oneof=1 2" msg:"冻结标识 1冻结 2解冻"`
}

type LoginHistoryItemResp struct {
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	Ip        string `json:"ip"`
	DeviceID  string `json:"device_id"`
	Platform  int64  ` json:"platform"`
	Brand     string `json:"brand"`
}

type LoginHistoryResp struct {
	List     []LoginHistoryItemResp `json:"list"`
	Count    int64                  `json:"count"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

type LoginHistoryReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id" form:"user_id" binding:"required" msg:"用户ID不能为空"`
	pagination.Pagination
}

type UserOnlineStatus struct {
	UserID string             `json:"user_id"`
	Type   common.MessageType `json:"type"`
	Msg    string             `json:"msg"`
}

type RealNameAuthReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	UserID      string `json:"user_id"`
	RealAuthMsg string `json:"real_auth_msg"`
	IsRealAuth  int    `json:"is_real_auth"`
}
