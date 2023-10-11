package model

import "im/pkg/pagination"

type RegistrationStatisticsReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	BeginDate   int64  `json:"begin_date" form:"start_date" binding:"omitempty,gte=6"  msg:"请输入开始正确的日期"`
	EndDate     int64  `json:"end_date" form:"end_date" binding:"omitempty,gte=6"  msg:"请输入结束正确的日期"`
	pagination.Pagination
}
type RegistrationStatisticsInfo struct {
	Daily string `json:"daily"`
	Count int64  `json:"count"`
}
type RegistrationStatisticsResp struct {
	List []RegistrationStatisticsInfo `json:"list"`
}

type InviteCodeStatisticsReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	InviteCode  string `json:"invite_code" form:"invite_code" binding:"omitempty,min=1"  msg:"请输入正确的邀请码"`
	pagination.Pagination
}

type InviteCodeStatisticsInfo struct {
	InviteCode string `json:"invite_code"`
	Count      int64  `json:"count"`
}

type InviteCodeStatisticsResp struct {
	List []InviteCodeStatisticsInfo `json:"list"`
	pagination.Pagination
}

type InviteCodeStatisticsDetailsReq struct {
	InviteCode  string `json:"invite_code" form:"invite_code" binding:"required,min=1,number"  msg:"请输入正确的邀请码"`
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	Account     string `json:"account" form:"account" binding:"omitempty,min=1"  msg:"请输入正确的帐号"`
	NickName    string `json:"nick_name" form:"nick_name" binding:"omitempty,min=4,max=16"  msg:"请输入正确的昵称"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"omitempty,min=10"  msg:"请输入正确的手机号码"`
	BeginDate   int64  `json:"begin_date" form:"begin_date" binding:"omitempty,min=1"  msg:"请输入正确的开始时间"`
	EndDate     int64  `json:"end_date" form:"end_date" binding:"omitempty,min=1"  msg:"请输入正确的结束时间"`
	pagination.Pagination
}

type InviteCodeStatisticsDetailsInfo struct {
	UserId          string  `json:"user_id"`
	Account         string  `json:"account"`
	PhoneNumber     string  `json:"phone_number"`
	FaceURL         string  `json:"face_url"`
	NickName        string  `json:"nick_name"`
	InviteCode      string  `json:"invite_code"`
	Balance         float64 `json:"balance"`
	RegistryTime    int64   `json:"registry_time"`
	LatestLoginTime int64   `json:"latest_login_time"`
}
type InviteCodeStatisticsDetailsResp struct {
	List []InviteCodeStatisticsDetailsInfo `json:"list"`
	pagination.Pagination
}

type SuggestionReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,min=1"`
	Account     string `json:"account" form:"account" binding:"omitempty,min=1"  msg:"请输入正确的帐号"`
	NickName    string `json:"nick_name" form:"nick_name" binding:"omitempty,min=1,max=16"  msg:"请输入正确的昵称"`
	UserId      string `json:"user_id" form:"user_id" binding:"omitempty,min=1,number"  msg:"请输入正确的user_id"`
	Content     string `json:"content" form:"content" binding:"omitempty,min=1"  msg:"请输入正确的内容"`
	Brand       string `json:"brand" form:"brand" binding:"omitempty,min=1"  msg:"请输入正确的品牌"`
	Platform    int64  `json:"platform" form:"platform" binding:"omitempty,min=1"  msg:"请选择正确的客户端类型"`
	BeginDate   int64  `json:"begin_date" form:"begin_date" binding:"omitempty,min=1"  msg:"请输入正确的开始时间"`
	EndDate     int64  `json:"end_date" form:"end_date" binding:"omitempty,min=1"  msg:"请输入正确的结束时间"`
	pagination.Pagination
}
type SuggestionInfo struct {
	Id         int64  `json:"id"`
	Account    string `json:"account"`
	NickName   string `json:"nick_name"`
	UserId     string `json:"user_id"`
	Content    string `json:"content"`
	Brand      string `json:"brand"`
	Platform   int64  `json:"platform"`
	AppVersion string `json:"app_version"`
	CreatedAt  int64  `json:"created_at"`
}
type SuggestionInfoResp struct {
	List []SuggestionInfo `json:"list"`
	pagination.Pagination
}

type OnlineUserInfo struct {
	UserID        string `json:"user_id"`
	Nickname      string `json:"nickname"`
	LastLoginTime int64  `json:"last_login_time"`
}

type OnlineUsersReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,min=1"`
	Keyword     string `json:"keyword" form:"keyword"`
	pagination.Pagination
}
type OnlineUsersResp struct {
	pagination.Pagination
	TotalCount int64            `json:"total_count"`
	List       []OnlineUserInfo `json:"list"`
}
